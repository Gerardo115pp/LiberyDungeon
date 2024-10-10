package handlers

import (
	"encoding/json"
	"fmt"
	"libery-dungeon-libs/communication"
	"libery-dungeon-libs/libs/libery_networking"
	dungeon_models "libery-dungeon-libs/models"
	app_config "libery_categories_service/Config"
	service_models "libery_categories_service/models"
	"libery_categories_service/repository"
	"libery_categories_service/workflows"
	common_workflows "libery_categories_service/workflows/common"
	"net/http"

	"github.com/Gerardo115pp/patriots_lib/echo"
)

func TrashcanHandler(service_instance libery_networking.Server) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		switch request.Method {
		case http.MethodGet:
			getTrashcanHandler(response, request)
		case http.MethodPost:
			postTrashcanHandler(response, request)
		case http.MethodPatch:
			patchTrashcanHandler(response, request)
		case http.MethodDelete:
			deleteTrashcanHandler(response, request)
		case http.MethodPut:
			putTrashcanHandler(response, request)
		case http.MethodOptions:
			response.WriteHeader(http.StatusOK)
		default:
			response.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func getTrashcanHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path

	switch resource {
	case "/trashcan/entries":
		getTrashcanEntriesHandler(response, request)
	case "/trashcan/transaction":
		getTrashcanTransactionHandler(response, request)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.getTrashcanHandler: resource<%s> not found", resource))
		response.WriteHeader(404)
	}
}

func getTrashcanEntriesHandler(response http.ResponseWriter, request *http.Request) {
	var transaction_entries []*service_models.TrashcanTransactionEntry = repository.TrashRepo.GetSortedTransactionEntries()
	var non_zero_entries []*service_models.TrashcanTransactionEntry = make([]*service_models.TrashcanTransactionEntry, 0)

	for _, entry := range transaction_entries {
		if entry.AffectedMedias > 0 {
			non_zero_entries = append(non_zero_entries, entry)
		}
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	json.NewEncoder(response).Encode(non_zero_entries)
}

func getTrashcanTransactionHandler(response http.ResponseWriter, request *http.Request) {
	var transaction_id string = request.URL.Query().Get("transaction_id")

	if transaction_id == "" {
		echo.Echo(echo.RedFG, "In handlers/trashcan.getTrashcanTransactionHandler: while getting transaction_id, it was empty")
		response.WriteHeader(400)
		return
	}

	transaction, err := repository.TrashRepo.GetTransaction(transaction_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.getTrashcanTransactionHandler: while getting transaction<%s>, error: %s", transaction_id, err.Error()))
		response.WriteHeader(404)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(200)
	json.NewEncoder(response).Encode(transaction)
}

func postTrashcanHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}

func patchTrashcanHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path

	switch resource {
	case "/trashcan/media/restore":
		patchTrashcanMediaRestoreHandler(response, request)
	case "/trashcan/transaction/restore":
		patchTrashcanTransactionRestoreHandler(response, request)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanHandler: resource<%s> not found", resource))
		response.WriteHeader(404)
	}
}

func patchTrashcanMediaRestoreHandler(response http.ResponseWriter, request *http.Request) {
	var media_uuid string = request.URL.Query().Get("media_uuid")
	var transaction_id string = request.URL.Query().Get("transaction_id")
	var new_main_category string = request.URL.Query().Get("main_category")
	var lerr *dungeon_models.LabeledError

	switch "" {
	case media_uuid:
		echo.Echo(echo.RedFG, "In handlers/trashcan.patchTrashcanMediaRestoreHandler: media_uuid is empty")
		response.WriteHeader(400)
		return
	case transaction_id:
		echo.Echo(echo.RedFG, "In handlers/trashcan.patchTrashcanMediaRestoreHandler: transaction_id is empty")
		response.WriteHeader(400)
		return
	case new_main_category:
		echo.Echo(echo.RedFG, "In handlers/trashcan.patchTrashcanMediaRestoreHandler: new_main_category is empty")
		response.WriteHeader(400)
		return
	}

	category_identity, lerr := common_workflows.GetCategoryIdentityFromUUID(new_main_category)
	if lerr != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanMediaRestoreHandler: while getting category_identity, error: %s", lerr.Error()))
		response.WriteHeader(404)
		return
	}

	transaction, err := repository.TrashRepo.GetTransaction(transaction_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanMediaRestoreHandler: while getting transaction<%s>, error: %s", transaction_id, err.Error()))
		response.WriteHeader(404)
		return
	}

	lerr = workflows.RestoreMediaFromTransaction(*category_identity, *transaction, media_uuid)
	if lerr != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanMediaRestoreHandler: while restoring media from transaction, error: %s", lerr.Error()))
		response.WriteHeader(500)
		return
	}

	repository.TrashRepo.EraseMediaFromTransaction(transaction_id, media_uuid)

	fs_change_event := communication.NewClusterFSChangeEvent(app_config.JWT_SECRET, category_identity.ClusterUUID, 0, 1, 0)
	err = fs_change_event.Emit()
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanMediaRestoreHandler: while emitting fs_change_event, error: %s", err.Error()))
	}

	response.WriteHeader(200)
}

func patchTrashcanTransactionRestoreHandler(response http.ResponseWriter, request *http.Request) {
	var transaction_id string = request.URL.Query().Get("transaction_id")
	var new_main_category string = request.URL.Query().Get("main_category")
	var lerr *dungeon_models.LabeledError

	if transaction_id == "" || new_main_category == "" {
		echo.Echo(echo.RedFG, "In handlers/trashcan.patchTrashcanTransactionRestoreHandler: transaction_id or new_main_category is empty")
		response.WriteHeader(400)
		return
	}

	category_identity, lerr := common_workflows.GetCategoryIdentityFromUUID(new_main_category)
	if lerr != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanTransactionRestoreHandler: while getting category_identity, error: %s", lerr.Error()))
		response.WriteHeader(404)
		return
	}

	transaction, err := repository.TrashRepo.GetTransaction(transaction_id)
	if err != nil {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanTransactionRestoreHandler: while getting transaction<%s>, error: %s", transaction_id, err.Error()))
		response.WriteHeader(404)
		return
	}

	go func() {
		medias_restored := len(transaction.Content)

		lerr = workflows.RestoreTransactionToCategory(*transaction, category_identity)
		if lerr != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanTransactionRestoreHandler: while restoring transaction to category, error: %s", lerr.Error()))
			response.WriteHeader(500)
			return
		}

		fs_change_event := communication.NewClusterFSChangeEvent(app_config.JWT_SECRET, category_identity.ClusterUUID, 0, medias_restored, 0)
		err = fs_change_event.Emit()
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.patchTrashcanTransactionRestoreHandler: while emitting fs_change_event, error: %s", err.Error()))
		}
	}()

	response.WriteHeader(200)
}

func deleteTrashcanHandler(response http.ResponseWriter, request *http.Request) {
	var resource string = request.URL.Path

	switch resource {
	case "/trashcan/transaction":
		deleteTrashcanTransactionHandler(response, request)
	case "/trashcan/empty":
		deleteEmptyTrashcanHandler(response, request)
	default:
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.deleteTrashcanHandler: resource<%s> not found", resource))
		response.WriteHeader(404)
	}
}

func deleteTrashcanTransactionHandler(response http.ResponseWriter, request *http.Request) {
	var transaction_id string = request.URL.Query().Get("transaction_id")

	if transaction_id == "" {
		echo.Echo(echo.RedFG, "In handlers/trashcan.deleteTrashcanTransactionHandler: while getting transaction_id, it was empty")
		response.WriteHeader(400)
		return
	}

	if !repository.TrashRepo.HasTransaction(transaction_id) {
		echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.deleteTrashcanTransactionHandler: transaction<%s> not found", transaction_id))
		response.WriteHeader(404)
		return
	}

	go func() {
		err := repository.TrashRepo.CleanSingleTransaction(transaction_id)
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.deleteTrashcanTransactionHandler: while deleting transaction<%s>, error: %s", transaction_id, err.Error()))
			response.WriteHeader(404)
			return
		}

		// Fs changes are not valid for this operation. we have to create a new event type for this.
	}()

	response.WriteHeader(200)
}

func deleteEmptyTrashcanHandler(response http.ResponseWriter, request *http.Request) {
	go func() {
		err := repository.TrashRepo.EmptyTrashcan()
		if err != nil {
			echo.Echo(echo.RedFG, fmt.Sprintf("In handlers/trashcan.deleteEmptyTrashcanHandler: error: %s", err.Error()))
			response.WriteHeader(500)
			return
		}
	}()

	response.WriteHeader(200)
}

func putTrashcanHandler(response http.ResponseWriter, request *http.Request) {
	response.WriteHeader(http.StatusMethodNotAllowed)
	return
}
