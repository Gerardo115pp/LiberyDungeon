#!/usr/bin/zsh

# This script is used to build the protobuf files for the project.
SHARED_DIR="/home/el_maligno/SoftwareProjects/LiberyDungeon/Services/Shared"
SHARED_GO_DIR="${SHARED_DIR}/go_shared"
SHARED_PROTO_DIR="${SHARED_DIR}/proto"
SHARED_PYTHON_DIR="${SHARED_DIR}/python_shared"

# INPUT SOURCES
PROTO_DOWNLOAD_SERVICE="${SHARED_PROTO_DIR}/download_service_pb"
PROTO_CATEGORIES_SERVICE="${SHARED_PROTO_DIR}/categories_service_pb"
PROTO_JD_SERVICE="${SHARED_PROTO_DIR}/jd_service_pb"
PROTO_METADATA_SERVICE="${SHARED_PROTO_DIR}/metadata_service_pb"

# OUTPUT SOURCES
PROTO_DOWNLOAD_SERVICE_GO="${SHARED_GO_DIR}/download_service_pb"
PROTO_CATEGORIES_SERVICE_GO="${SHARED_GO_DIR}/categories_service_pb"
PROTO_JD_SERVICE_GO="${SHARED_GO_DIR}/jd_service_pb"
PROTO_METADATA_SERVICE_GO="${SHARED_GO_DIR}/metadata_service_pb"
PROTO_DOWNLOAD_SERVICE_PYTHON="${SHARED_PYTHON_DIR}/DownloadServicePB"
PROTO_CATEGORIES_SERVICE_PYTHON="${SHARED_PYTHON_DIR}/CategoriesServicePB"
PROTO_JD_SERVICE_PYTHON="${SHARED_PYTHON_DIR}/JDServicePB"
PROTO_METADATA_SERVICE_PYTHON="${SHARED_PYTHON_DIR}/MetadataServicePB"


# CHECK IF THE OUTPUT DIRECTORIES EXIST
if [ ! -d "${PROTO_DOWNLOAD_SERVICE_GO}" ]; then
    mkdir -p "${PROTO_DOWNLOAD_SERVICE_GO}"
fi

if [ ! -d "${PROTO_CATEGORIES_SERVICE_GO}" ]; then
    mkdir -p "${PROTO_CATEGORIES_SERVICE_GO}"
fi

if [ ! -d "${PROTO_JD_SERVICE_GO}" ]; then
    mkdir -p "${PROTO_JD_SERVICE_GO}"
fi

if [ ! -d "${PROTO_METADATA_SERVICE_GO}" ]; then
    mkdir -p "${PROTO_METADATA_SERVICE_GO}"
fi

if [ ! -d "${PROTO_DOWNLOAD_SERVICE_PYTHON}" ]; then
    mkdir -p "${PROTO_DOWNLOAD_SERVICE_PYTHON}"
    touch "${PROTO_DOWNLOAD_SERVICE_PYTHON}/__init__.py"
fi

if [ ! -d "${PROTO_CATEGORIES_SERVICE_PYTHON}" ]; then
    mkdir -p "${PROTO_CATEGORIES_SERVICE_PYTHON}"
    touch "${PROTO_CATEGORIES_SERVICE_PYTHON}/__init__.py"
fi

if [ ! -d "${PROTO_JD_SERVICE_PYTHON}" ]; then
    mkdir -p "${PROTO_JD_SERVICE_PYTHON}"
    touch "${PROTO_JD_SERVICE_PYTHON}/__init__.py"
fi

# GENERATE THE PROTO FILES FOR GO
generate_go_protos() {
    protoc --proto_path="${PROTO_DOWNLOAD_SERVICE}" --go_out="${PROTO_DOWNLOAD_SERVICE_GO}" --go_opt=paths=source_relative --go-grpc_out="${PROTO_DOWNLOAD_SERVICE_GO}" --go-grpc_opt=paths=source_relative "${PROTO_DOWNLOAD_SERVICE}/download_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the download service on go"   
        exit 1
    fi

    protoc --proto_path="${PROTO_CATEGORIES_SERVICE}" --go_out="${PROTO_CATEGORIES_SERVICE_GO}" --go_opt=paths=source_relative --go-grpc_out="${PROTO_CATEGORIES_SERVICE_GO}" --go-grpc_opt=paths=source_relative "${PROTO_CATEGORIES_SERVICE}/categories_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the categories service on go"   
        exit 1
    fi

    protoc --proto_path="${PROTO_JD_SERVICE}" --go_out="${PROTO_JD_SERVICE_GO}" --go_opt=paths=source_relative --go-grpc_out="${PROTO_JD_SERVICE_GO}" --go-grpc_opt=paths=source_relative "${PROTO_JD_SERVICE}/jd_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the jd service on go"   
        exit 1
    fi

    protoc --proto_path="${PROTO_METADATA_SERVICE}" --go_out="${PROTO_METADATA_SERVICE_GO}" --go_opt=paths=source_relative --go-grpc_out="${PROTO_METADATA_SERVICE_GO}" --go-grpc_opt=paths=source_relative "${PROTO_METADATA_SERVICE}/metadata_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the metadata service on go"   
        exit 1
    fi
}

generate_python_protos() {
    python -m grpc_tools.protoc -I "${PROTO_DOWNLOAD_SERVICE}" --python_out="${PROTO_DOWNLOAD_SERVICE_PYTHON}" --grpc_python_out="${PROTO_DOWNLOAD_SERVICE_PYTHON}" --pyi_out="${PROTO_DOWNLOAD_SERVICE_PYTHON}" "${PROTO_DOWNLOAD_SERVICE}/download_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the download service on python"   
        exit 1
    fi
    # Change the relative import to absolute import
    # so 'import download_requests_pb2 as download__requests__pb2' -> 'import DownloadServicePB.download_requests_pb2 as download__requests__pb2'
    sed -i 's/import download_requests_pb2 as download__requests__pb2/import DownloadServicePB.download_requests_pb2 as download__requests__pb2/g' "${PROTO_DOWNLOAD_SERVICE_PYTHON}/download_requests_pb2_grpc.py"

    python -m grpc_tools.protoc -I "${PROTO_CATEGORIES_SERVICE}" --python_out="${PROTO_CATEGORIES_SERVICE_PYTHON}" --grpc_python_out="${PROTO_CATEGORIES_SERVICE_PYTHON}" --pyi_out="${PROTO_CATEGORIES_SERVICE_PYTHON}" "${PROTO_CATEGORIES_SERVICE}/categories_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the categories service on python"   
        exit 1
    fi
    # Change the relative import to absolute import
    # so 'import categories_requests_pb2 as categories__requests__pb2' -> 'import CategoriesServicePB.categories_requests_pb2 as categories__requests__pb2'
    sed -i 's/import categories_requests_pb2 as categories__requests__pb2/import CategoriesServicePB.categories_requests_pb2 as categories__requests__pb2/g' "${PROTO_CATEGORIES_SERVICE_PYTHON}/categories_requests_pb2_grpc.py"

    python -m grpc_tools.protoc -I "${PROTO_JD_SERVICE}" --python_out="${PROTO_JD_SERVICE_PYTHON}" --grpc_python_out="${PROTO_JD_SERVICE_PYTHON}" --pyi_out="${PROTO_JD_SERVICE_PYTHON}" "${PROTO_JD_SERVICE}/jd_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the jd service on python"   
        exit 1
    fi
    # Change the relative import to absolute import
    # so 'import jd_requests_pb2 as jd__requests__pb2' -> 'import JDServicePB.jd_requests_pb2 as jd__requests__pb2'
    sed -i 's/import jd_requests_pb2 as jd__requests__pb2/import JDServicePB.jd_requests_pb2 as jd__requests__pb2/g' "${PROTO_JD_SERVICE_PYTHON}/jd_requests_pb2_grpc.py"

    python -m grpc_tools.protoc -I "${PROTO_METADATA_SERVICE}" --python_out="${PROTO_METADATA_SERVICE_PYTHON}" --grpc_python_out="${PROTO_METADATA_SERVICE_PYTHON}" --pyi_out="${PROTO_METADATA_SERVICE_PYTHON}" "${PROTO_METADATA_SERVICE}/metadata_requests.proto"
    if [ $? -ne 0 ]; then
        echo "Error compiling the proto files for the metadata service on python"   
        exit 1
    fi
    # Change the relative import to absolute import
    # so 'import metadata_requests_pb2 as metadata__requests__pb2' -> 'import MetadataServicePB.metadata_requests_pb2 as metadata__requests__pb2'
    sed -i 's/import metadata_requests_pb2 as metadata__requests__pb2/import MetadataServicePB.metadata_requests_pb2 as metadata__requests__pb2/g' "${PROTO_METADATA_SERVICE_PYTHON}/metadata_requests_pb2_grpc.py"
}

generate_go_protos
generate_python_protos