package cli

import (
	"fmt"
	"os"
)

// runCompletion handles the "ath completion" command.
// Currently only zsh is supported.
func runCompletion(args []string) int {
	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Usage: ath completion zsh\n")
		return 2
	}

	switch args[0] {
	case "zsh":
		fmt.Print(zshCompletionScript)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unsupported shell: %s (only zsh is supported)\n", args[0])
		return 2
	}
}

const zshCompletionScript = `#compdef ath

# Zsh completion for ath — athanor agent orchestration CLI
# Install: ath completion zsh > ~/.zsh/completions/_ath

_ath_athanor_names() {
    local -a names
    names=( ${(f)"$(ls ~/athanor/athanors/ 2>/dev/null)"} )
    compadd -- "${names[@]}"
}

_ath_opus_files() {
    local athanor_dir
    athanor_dir="${ATHANOR:-}"
    if [[ -z "$athanor_dir" ]]; then
        return
    fi
    local -a files
    files=( ${(f)"$(ls "$athanor_dir/opera/" 2>/dev/null)"} )
    compadd -- "${files[@]}"
}

_ath_tmux_windows() {
    local -a windows
    windows=( ${(f)"$(tmux list-windows -a -F '#{window_name}' 2>/dev/null)"} )
    compadd -- "${windows[@]}"
}

_ath_tmux_azer_windows() {
    local -a windows
    windows=( ${(f)"$(tmux list-windows -a -F '#{window_name}' 2>/dev/null | grep '^azer-')"} )
    compadd -- "${windows[@]}"
}

_ath_mo_names() {
    local athanor_name="$1"
    local -a names
    # Multi-MO: list files from magna-opera/
    names=( ${(f)"$(ls ~/athanor/athanors/$athanor_name/magna-opera/*.md 2>/dev/null | xargs -I{} basename {} .md)"} )
    if [[ ${#names[@]} -eq 0 ]]; then
        # Legacy: no completion needed (mo-name is optional)
        return
    fi
    compadd -- "${names[@]}"
}

_ath() {
    local -a commands
    commands=(
        'init:Create a new athanor instance'
        'kindle:Launch a marut for an athanor'
        'reforge:Kill and relaunch a marut'
        'muster:Launch an azer for an opus'
        'cleanup:Clean up after a discharged opus'
        'quiesce:Graceful shutdown of an athanor'
        'status:Show athanor health'
        'opera:List opera with status'
        'whisper:Reliable message delivery to tmux sessions'
        'completion:Generate shell completion script'
        'version:Print version info'
    )

    if (( CURRENT == 2 )); then
        _describe -t commands 'ath command' commands
        return
    fi

    case "${words[2]}" in
        kindle|reforge|quiesce)
            if (( CURRENT == 3 )); then
                _ath_athanor_names
            elif (( CURRENT == 4 )); then
                _ath_mo_names "${words[3]}"
            fi
            ;;
        status|opera)
            if (( CURRENT == 3 )); then
                _ath_athanor_names
            fi
            ;;
        init)
            # init takes a name (free text), then optional flags
            if (( CURRENT == 3 )); then
                _message 'athanor name'
            elif (( CURRENT >= 4 )); then
                _arguments '*--project[Working directory]:directory:_directories'
            fi
            ;;
        muster)
            if (( CURRENT == 3 )); then
                _ath_opus_files
            fi
            ;;
        cleanup)
            if (( CURRENT == 3 )); then
                _ath_tmux_azer_windows
            fi
            ;;
        whisper)
            local -a whisper_commands
            whisper_commands=(
                'send:Send a message to a tmux target'
                'idle:Wait for target to become idle'
                'wait-and-send:Wait for idle, then send'
            )
            if (( CURRENT == 3 )); then
                _describe -t commands 'whisper command' whisper_commands
            elif (( CURRENT == 4 )); then
                case "${words[3]}" in
                    send|idle|wait-and-send)
                        _ath_tmux_windows
                        ;;
                esac
            fi
            ;;
        completion)
            if (( CURRENT == 3 )); then
                compadd -- zsh
            fi
            ;;
    esac
}

_ath "$@"
`
