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
		os.Stdout.WriteString(zshCompletionScript)
		return 0
	default:
		fmt.Fprintf(os.Stderr, "unsupported shell: %s (only zsh is supported)\n", args[0])
		return 2
	}
}

const zshCompletionScript = `#compdef ath

# Zsh completion for ath — athanor agent orchestration CLI
# Install: ath completion zsh > ~/.zsh/completions/_ath

# Ensure _message hints are visible (e.g. "session name" for free-text args)
zstyle ':completion:*:ath:*:messages' format '%d'

_ath_athanor_names() {
    local -a names
    names=( ${(f)"$(ls ~/athanor/athanors/ 2>/dev/null)"} )
    compadd -- "${names[@]}"
}

_ath_charged_opus_files() {
    local athanor_dir="$1"
    local -a charged charged_d
    local opera_dir mo_name f name opus_status
    for opera_dir in "$athanor_dir"/magna-opera/*/opera(N/); do
        mo_name="${${opera_dir:h}:t}"
        for f in "$opera_dir"/*.md(N); do
            name="${f:t}"
            opus_status=$(awk '/^---$/{if(n++)exit}n&&/^status:/{print $2}' "$f")
            if [[ "$opus_status" == "charged" ]]; then
                charged+=("$name")
                charged_d+=("[$mo_name] ${name%.md}")
            fi
        done
    done
    (( ${#charged} )) && compadd -V charged -d charged_d -- "${charged[@]}"
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
    # Multi-MO: list directories from magna-opera/
    names=( ${(f)"$(ls -d ~/athanor/athanors/$athanor_name/magna-opera/*/ 2>/dev/null | xargs -I{} basename {})"} )
    if [[ ${#names[@]} -eq 0 ]]; then
        # Legacy: no completion needed (mo-name is optional)
        return
    fi
    compadd -- "${names[@]}"
}

_ath_opus_names_ordered() {
    local athanor_name="$1" mo_name="$2"
    local opera_dir="$HOME/athanor/athanors/$athanor_name/magna-opera/$mo_name/opera"
    local -a charged charged_d discharged discharged_d assessed assessed_d other other_d

    if [[ ! -d "$opera_dir" ]]; then
        return
    fi

    for f in "$opera_dir"/*.md(N); do
        local name="${${f:t}%.md}"
        local opus_status=$(awk '/^---$/{if(n++)exit}n&&/^status:/{print $2}' "$f")
        case "$opus_status" in
            charged) charged+=("$name"); charged_d+=("[charged] $name") ;;
            discharged) discharged+=("$name"); discharged_d+=("[discharged] $name") ;;
            assessed) assessed+=("$name"); assessed_d+=("[assessed] $name") ;;
            *) other+=("$name"); other_d+=("[?] $name") ;;
        esac
    done

    [[ ${#charged} -gt 0 ]] && compadd -V charged -d charged_d -- "${charged[@]}"
    [[ ${#discharged} -gt 0 ]] && compadd -V discharged -d discharged_d -- "${discharged[@]}"
    [[ ${#assessed} -gt 0 ]] && compadd -V assessed -d assessed_d -- "${assessed[@]}"
    [[ ${#other} -gt 0 ]] && compadd -V other -d other_d -- "${other[@]}"
}

_ath() {
    local -a commands
    commands=(
        'init:Create a new athanor instance'
        'craft:Interactive session with the artifex'
        'craft-mo:Create a new Magnum Opus interactively'
        'kindle:Launch a marut for an athanor'
        'reforge:Kill and relaunch a marut'
        'muster:Launch an azer for an opus'
        'check:Check crucible health'
        'cleanup:Clean up after a discharged opus'
        'quiesce:Graceful shutdown of an athanor'
        'status:Show athanor health'
        'view:Open MO or opus in \$EDITOR'
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
        craft)
            if (( CURRENT == 3 )); then
                _ath_athanor_names
            elif (( CURRENT == 4 )); then
                _ath_mo_names "${words[3]}"
            elif (( CURRENT == 5 )); then
                _message 'session name'
            fi
            ;;
        craft-mo)
            if (( CURRENT == 3 )); then
                _ath_athanor_names
            fi
            ;;
        kindle|reforge|quiesce)
            if (( CURRENT == 3 )); then
                _ath_athanor_names
            elif (( CURRENT == 4 )); then
                _ath_mo_names "${words[3]}"
            fi
            ;;
        view)
            if (( CURRENT == 3 )); then
                _ath_athanor_names
            elif (( CURRENT == 4 )); then
                _ath_mo_names "${words[3]}"
            elif (( CURRENT == 5 )); then
                _ath_opus_names_ordered "${words[3]}" "${words[4]}"
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
            # Resolve athanor dir from --athanor flag or $ATHANOR
            local _ath_muster_dir="${ATHANOR:-}"
            local -i _ath_flag_idx=${words[(I)--athanor]}
            if (( _ath_flag_idx && _ath_flag_idx + 1 < CURRENT )); then
                _ath_muster_dir="$HOME/athanor/athanors/${words[$(( _ath_flag_idx + 1 ))]}"
            fi

            # Complete flag values
            case "${words[$(( CURRENT - 1 ))]}" in
                --athanor) _ath_athanor_names; return ;;
                --dir) _directories; return ;;
                --model) compadd -- sonnet opus haiku; return ;;
                --name|--intent) return ;;
            esac

            # Count positional args (skip flags and their values)
            local -i _pos=0 _i _skip=0
            for (( _i = 3; _i < CURRENT; _i++ )); do
                if (( _skip )); then _skip=0; continue; fi
                case "${words[$_i]}" in
                    --athanor|--dir|--model|--name|--intent) _skip=1 ;;
                    -*) ;;
                    *) (( _pos++ )) ;;
                esac
            done

            # Collect available flags (exclude already-used ones)
            local -a _mflags=(--athanor --dir --model --name --intent)
            for (( _i = 3; _i < CURRENT; _i++ )); do
                _mflags=("${(@)_mflags:#${words[$_i]}}")
            done

            # Positional completions depend on mode
            if (( ${words[(I)--intent]} )); then
                # Intent mode: <mo> <name> --intent <text>
                case $_pos in
                    0)
                        if [[ -n "$_ath_muster_dir" && -d "$_ath_muster_dir" ]]; then
                            _ath_mo_names "${_ath_muster_dir:t}"
                        else
                            _message 'MO name (use --athanor <name> first)'
                        fi
                        ;;
                    1) _message 'session name' ;;
                esac
            else
                # Opus mode: <opus-file>
                case $_pos in
                    0)
                        if [[ -n "$_ath_muster_dir" && -d "$_ath_muster_dir" ]]; then
                            _ath_charged_opus_files "$_ath_muster_dir"
                        else
                            _message 'opus file (use --athanor <name> or set $ATHANOR)'
                        fi
                        ;;
                esac
            fi
            # Always offer remaining flags
            (( ${#_mflags} )) && compadd -- "${_mflags[@]}"
            ;;
        check)
            if (( CURRENT == 3 )); then
                _ath_tmux_windows
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
