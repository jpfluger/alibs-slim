#!/bin/bash

# Always run this script from inside the `maintainers` directory.

# Example usage:
# 1. Dry-run full workflow (debug mode)
#    ./versionapply.sh -D -U -T -P -G "bump release" -B main
# 2. Full workflow (update, tag locally, push tags, commit changes)
#    ./versionapply.sh -U -T -P -G "bump release" -B main
# 3. Update versions and tag locally (without pushing)
#    ./versionapply.sh -U -T
# 4. Push previously created tags to remote (no new updates or commits)
#    ./versionapply.sh -P
# 5. Only commit and push existing changes (no updates or tags)
#    ./versionapply.sh -G "Fix dependencies" -B dev

# Initialize variables
DEBUG_MODE=false
TAG_MODE=false
UPDATE_MODE=false
GIT_COMMIT_MODE=false
GIT_COMMIT_MESSAGE=""
PUSH_TAGS_MODE=false
BRANCH_NAME="main"

usage() {
    echo "usage: $0 [options]"
    echo "Options:"
    echo "  -h, --help            Show help"
    echo "  -D, --debug           Debug mode (no changes applied)"
    echo "  -U, --update-version  Update version numbers in go.mod files"
    echo "  -T, --tag             Create Git tags locally"
    echo "  -P, --push-tags       Push Git tags to remote"
    echo "  -G, --git [msg]       Commit and push changes with provided message"
    echo "  -B, --branch [name]   Specify branch name (default: main)"
}

# Handle parameters
while [[ "$#" -gt 0 ]]; do
    case $1 in
        -D|--debug) DEBUG_MODE=true ;;
        -U|--update-version) UPDATE_MODE=true ;;
        -T|--tag) TAG_MODE=true ;;
        -P|--push-tags) PUSH_TAGS_MODE=true ;;
        -G|--git)
            GIT_COMMIT_MODE=true; shift
            [[ -n "$1" ]] && GIT_COMMIT_MESSAGE="$1" || { echo "ERROR: --git requires a commit message."; exit 1; } ;;
        -B|--branch)
            shift; [[ -n "$1" ]] && BRANCH_NAME="$1" || { echo "ERROR: --branch requires a branch name."; exit 1; } ;;
        -h|--help) usage; exit 0 ;;
        *) echo "Invalid option: $1"; usage; exit 1 ;;
    esac
    shift
done

# Load version from ./version file
if [ ! -f ./version ]; then
    echo "Error: ./version file not found."
    exit 1
fi

VERSION=$(cat ./version)
if [ -z "$VERSION" ]; then
    echo "Error: version number is empty."
    exit 1
fi

echo "Loaded version: $VERSION"
tag_name="v$VERSION"

# Move to parent directory
cd ..

# Pre-check tags function
check_existing_tags() {
    local tag="$1"
    if git rev-parse "$tag" >/dev/null 2>&1; then
        echo "ERROR: Tag $tag already exists in root repo. Exiting."
        exit 1
    fi
    find aconns -name "go.mod" | while read -r go_mod_file; do
        local mod_dir=$(dirname "$go_mod_file")
        local rel_path=$(realpath --relative-to="." "$mod_dir")
        local sub_tag="$rel_path/$tag"
        if git -C "$mod_dir" rev-parse "$sub_tag" >/dev/null 2>&1; then
            echo "ERROR: Tag $sub_tag exists in $mod_dir. Exiting."
            exit 1
        fi
    done
    echo "No existing tags found for $tag. Proceeding."
}

# Tag creation functions
create_git_tag() {
    local dir=$1
    local version=$2
    local rel_path=$(realpath --relative-to="." "$dir")

    # Special case: root directory (.)
    if [ "$rel_path" = "." ]; then
        tag="v$version"
    else
        tag="$rel_path/v$version"
    fi

    if git -C "$dir" rev-parse "$tag" >/dev/null 2>&1; then
        echo "Tag $tag already exists in $dir. Skipping."
        return
    fi

    echo "Creating Git tag: $tag"
    if ! $DEBUG_MODE; then
        git -C "$dir" tag "$tag"
        if $PUSH_TAGS_MODE; then
            git -C "$dir" push origin "$tag" || { echo "ERROR: Push failed for $tag in $dir"; exit 1; }
        fi
    else
        echo "[DEBUG] Would create tag $tag in $dir"
    fi
}

create_root_tag() {
    local version=$1
    local tag="v$version"

    if git rev-parse "$tag" >/dev/null 2>&1; then
        echo "Root tag $tag already exists. Skipping."
        return
    fi

    echo "Creating Git root tag: $tag"
    if ! $DEBUG_MODE; then
        git tag "$tag"
        if $PUSH_TAGS_MODE; then
            git push origin "$tag" || { echo "ERROR: Push failed for root tag $tag"; exit 1; }
        fi
    else
        echo "[DEBUG] Would create root tag $tag"
    fi
}

# Update go.mod files
if $UPDATE_MODE; then
    find aconns -name "go.mod" | while read -r go_mod_file; do
        echo "Processing $go_mod_file"
        MOD_DIR=$(dirname "$go_mod_file")
        TMP_FILE="$MOD_DIR/go.mod.tmp"
        inside_require=false

        while IFS= read -r line; do
            [[ "$line" =~ ^require\ \( ]] && inside_require=true
            [[ "$line" =~ ^\) ]] && inside_require=false

            if $inside_require && [[ "$line" =~ github.com/jpfluger/alibs-slim ]]; then
                module=$(echo "$line" | awk '{print $1}')
                indent=$(echo "$line" | sed 's/\S.*//')
                new_line="$indent$module v$VERSION"
                if ! $DEBUG_MODE; then
                    echo "$new_line" >> "$TMP_FILE"
                else
                    echo "[DEBUG] Would update $line to $new_line"
                fi
            else
                if ! $DEBUG_MODE; then
                    echo "$line" >> "$TMP_FILE"
                fi
            fi
        done < "$go_mod_file"

        if ! $DEBUG_MODE; then
            mv "$TMP_FILE" "$go_mod_file"
            pushd "$MOD_DIR" > /dev/null
            go mod tidy
            popd > /dev/null
        fi
    done
    echo "All go.mod files updated."
fi

# Git commit and push changes
# We are committing at the root. We aren't using `submodules`.
# Our setup is a single repository containing multiple Go modules.
if $GIT_COMMIT_MODE; then
    echo "Committing changes with message: '$GIT_COMMIT_MESSAGE'"
    if ! $DEBUG_MODE; then
        git add -A
        git commit -m "$GIT_COMMIT_MESSAGE"
        git push origin "$BRANCH_NAME" || { echo "ERROR: Failed to push changes."; exit 1; }
    else
        echo "[DEBUG] Would commit '$GIT_COMMIT_MESSAGE' and push to $BRANCH_NAME"
    fi
fi

# Tagging actions
if $TAG_MODE; then
    check_existing_tags "$tag_name"
    create_root_tag "$VERSION"
    find aconns -name "go.mod" | while read -r go_mod_file; do
        MOD_DIR=$(dirname "$go_mod_file")
        create_git_tag "$MOD_DIR" "$VERSION"
    done
    echo "All tags created successfully."
fi

echo "Script execution completed."
