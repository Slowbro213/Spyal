#!/bin/bash

set -euo pipefail

error_exit() {
    echo "❌ $1" >&2
    exit 1
}

success_msg() {
    echo "✅ $1"
}

ENV_FILE="../.env.public"
if [[ -f "$ENV_FILE" ]]; then
    echo "📄 Loading environment variables from $ENV_FILE"
    # Export all non-comment variables
    set -o allexport
    source "$ENV_FILE"
    set +o allexport
    success_msg "Environment variables loaded"
else
    echo "⚠️ $ENV_FILE not found"
    exit 1
fi

if ! command -v bun &> /dev/null; then
    error_exit "bun is not installed or not in PATH. Please install bun first."
fi


shopt -s nullglob
scripts=(scripts/build/*.{ts,sh})
shopt -u nullglob


if [[ ${#scripts[@]} -eq 0 ]]; then
    error_exit "No TypeScript files found in scripts/build/"
fi

echo "🔨 Running build scripts..."


for script in "${scripts[@]}"; do
    [[ $(basename "$script") == "build.sh" ]] && continue
    scripts+=("$script")
    if [[ -f "$script" ]]; then
        echo "📦 Running $script..."
        script_name=$(basename "$script")

        case "$script" in
            *.ts)
                if ! bun run "$script"; then
                    error_exit "Script $script_name failed with exit code $?"
                fi
                ;;
            *.sh)
                if ! bash "$script"; then
                    error_exit "Script $script_name failed with exit code $?"
                fi
                ;;
            *)
                echo "⚠️ Skipping unknown script type: $script_name"
                ;;
        esac

        success_msg "$script_name completed successfully"
    else
        echo "⚠️ Warning: $script is not a valid file, skipping"
    fi
done

success_msg "All build scripts completed successfully! 🎉"
