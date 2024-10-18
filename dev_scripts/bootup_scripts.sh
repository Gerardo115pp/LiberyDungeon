DUNGEON_ROOT_DIR=/home/el_maligno/SoftwareProjects/LiberyDungeon

function run_frontend {
    local frontend_dir=$DUNGEON_ROOT_DIR/apps/libery_dungeon

    if [ ! -d $frontend_dir ]; then
        echo "Frontend directory not found: $frontend_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env

    dotenv-cli -e $global_env_file -- npm run --prefix $frontend_dir dev
}

function run_medias {
    local medias_dir=$DUNGEON_ROOT_DIR/Services/Medias

    if [ ! -d $medias_dir ]; then
        echo "Medias directory not found: $medias_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env
    local medias_env_file=$env_directory/secretfile_medias.env

    local media_service_main=$medias_dir/medias_service.go

    if [ ! -f $media_service_main ]; then
        echo "Media service main file not found: $media_service_main"
        return 1
    fi

    local current_pwd=$(pwd)

    cd $medias_dir

    dotenv-cli -e $global_env_file -e $medias_env_file -- go run $media_service_main

    cd $current_pwd
}

function run_categories {
    local categories_dir=$DUNGEON_ROOT_DIR/Services/Categories

    if [ ! -d $categories_dir ]; then
        echo "Categories directory not found: $categories_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env
    local categories_env_file=$env_directory/secretfile_categories.env

    local categories_service_main=$categories_dir/categories_service.go

    if [ ! -f $categories_service_main ]; then
        echo "Categories service main file not found: $categories_service_main"
        return 1
    fi

    local current_pwd=$(pwd)

    cd $categories_dir

    dotenv-cli -e $global_env_file -e $categories_env_file -- go run $categories_service_main

    cd $current_pwd
}

function run_JD {
    local JD_dir=$DUNGEON_ROOT_DIR/Services/JD

    if [ ! -d $JD_dir ]; then
        echo "JD directory not found: $JD_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env
    local JD_env_file=$env_directory/secretfile_JD.env

    local JD_service_main=$JD_dir/JD_service.go

    if [ ! -f $JD_service_main ]; then
        echo "JD service main file not found: $JD_service_main"
        return 1
    fi

    local current_pwd=$(pwd)

    cd $JD_dir

    dotenv-cli -e $global_env_file -e $JD_env_file -- go run $JD_service_main

    cd $current_pwd
}

function run_metadata {
    local metadata_dir=$DUNGEON_ROOT_DIR/Services/Metadata

    if [ ! -d $metadata_dir ]; then
        echo "Metadata directory not found: $metadata_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env
    local metadata_env_file=$env_directory/secretfile_metadata.env

    local metadata_service_main=$metadata_dir/metadata_service.go

    if [ ! -f $metadata_service_main ]; then
        echo "Metadata service main file not found: $metadata_service_main"
        return 1
    fi

    local current_pwd=$(pwd)

    cd $metadata_dir

    dotenv-cli -e $global_env_file -e $metadata_env_file -- go run $metadata_service_main

    cd $current_pwd
}

function run_users {
    local users_dir=$DUNGEON_ROOT_DIR/Services/Users

    if [ ! -d $users_dir ]; then
        echo "Users directory not found: $users_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env
    local users_env_file=$env_directory/secretfile_users.env

    local users_service_main=$users_dir/users_service.go

    if [ ! -f $users_service_main ]; then
        echo "Users service main file not found: $users_service_main"
        return 1
    fi

    local current_pwd=$(pwd)

    cd $users_dir

    dotenv-cli -e $global_env_file -e $users_env_file -- go run $users_service_main

    cd $current_pwd
}

function run_downloads {
    local downloads_dir=$DUNGEON_ROOT_DIR/Services/Downloads

    if [ ! -d $downloads_dir ]; then
        echo "Downloads directory not found: $downloads_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env
    local downloads_env_file=$env_directory/secretfile_downloads.env

    local downloads_service_main=$downloads_dir/downloads_service.go

    if [ ! -f $downloads_service_main ]; then
        echo "Downloads service main file not found: $downloads_service_main"
        return 1
    fi

    local current_pwd=$(pwd)

    cd $downloads_dir

    dotenv-cli -e $global_env_file -e $downloads_env_file -- go run $downloads_service_main

    cd $current_pwd
}

function run_collect {
    local collect_dir=$DUNGEON_ROOT_DIR/Services/Collect

    if [ ! -d $collect_dir ]; then
        echo "Collect directory not found: $collect_dir"
        return 1
    fi

    local env_directory=$DUNGEON_ROOT_DIR/dev_scripts/env

    if [ ! -d $env_directory ]; then
        echo "Environment directory not found: $env_directory"
        return 1
    fi

    local global_env_file=$env_directory/secretfile_global.env
    local collect_env_file=$env_directory/secretfile_collect.env

    local collect_service_main=$collect_dir/collect_service.py

    if [ ! -f $collect_service_main ]; then
        echo "Collect service main file not found: $collect_service_main"
        return 1
    fi

    local current_pwd=$(pwd)

    cd $collect_dir

    local DEVELOPMENT_PYTHONPATH="{$PYTHONPATH}:$DUNGEON_ROOT_DIR/Services/Shared/python_shared"

    dotenv-cli -e $global_env_file -e $collect_env_file -v PYTHONPATH=$DEVELOPMENT_PYTHONPATH -- uvicorn collect_service:app --reload --port 6972 --host 0.0.0.0

    cd $current_pwd
}

function run_services {
    # JD has to be run first so there is no point on adding it here.
    run_medias &
    run_categories &
    run_metadata &
    run_users &
    run_downloads &
    run_collect
}