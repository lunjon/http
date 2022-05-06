_http_complete() {
    local nargs=${#COMP_WORDS[@]}

    local root_commands="get post put delete head options patch version alias"
    if [[ $nargs -eq 1 ]]; then
        COMPREPLY+=($(compgen -W "$root_commands"))
        return
    fi

    local cur=${COMP_WORDS[COMP_CWORD]}
    COMPREPLY=()

    case "$cur" in
        version|alias)
            return
            ;;
        -*)
            COMPREPLY+=($(compgen -W '-v --verbose -f --fail -s --silent \
                                      -T --timeout --key --cert' -- $cur))
            ;;
        *)
            COMPREPLY+=($(compgen -W '-v --verbose -f --fail -s --silent \
                                      -T --timeout --key --cert' -- $cur))

            has_url=0
            for arg in "${COMP_LINE[@]:1}"; do
                if [[ "$arg" =~ ^https?://[a-z]+$ ]]; then
                    has_url=1
                    break
                fi
            done

            if [[ $has_url -eq 0 ]]; then
                local urls=$(http alias | awk '{ print $2 }' | xargs)
                COMPREPLY+=($(compgen -W "$urls"))
            fi
            ;;
    esac
}

complete -F _http_complete http
