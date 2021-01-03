package cli

const (
	bashCompletionFunc = `# call arvan paas get $1,
__arvan_paas_override_flag_list=(config cluster user context namespace server)
__arvan_paas_override_flags()
{
    local ${__arvan_paas_override_flag_list[*]} two_word_of of
    for w in "${words[@]}"; do
        if [ -n "${two_word_of}" ]; then
            eval "${two_word_of}=\"--${two_word_of}=\${w}\""
            two_word_of=
            continue
        fi
        for of in "${__arvan_paas_override_flag_list[@]}"; do
            case "${w}" in
                --${of}=*)
                    eval "${of}=\"${w}\""
                    ;;
                --${of})
                    two_word_of="${of}"
                    ;;
            esac
        done
    done
    for of in "${__arvan_paas_override_flag_list[@]}"; do
        if eval "test -n \"\$${of}\""; then
            eval "echo \${${of}}"
        fi
    done
}
__arvan_paas_parse_get()
{

    local template
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    local arvan_paas_out
    if arvan_paas_out=$(arvan paas get $(__arvan_paas_override_flags) -o template --template="${template}" "$1" 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${arvan_paas_out[*]}" -- "$cur" ) )
    fi
}

__arvan_paas_get_namespaces()
{
    local template arvan_paas_out
    template="{{ range .items  }}{{ .metadata.name }} {{ end }}"
    if arvan_paas_out=$(arvan paas get -o template --template="${template}" namespace 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${arvan_paas_out[*]}" -- "$cur" ) )
    fi
}

__arvan_paas_get_resource()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
      local arvan_paas_out
      if arvan_paas_out=$(arvan paas api-resources $(__arvan_paas_override_flags) -o name --cached --request-timeout=5s --verbs=get 2>/dev/null); then
          COMPREPLY=( $( compgen -W "${arvan_paas_out[*]}" -- "$cur" ) )
          return 0
      fi
      return 1
    fi
    __arvan_paas_parse_get "${nouns[${#nouns[@]} -1]}"
}

# $1 is the name of the pod we want to get the list of containers inside
__arvan_paas_get_containers()
{
    local template
    template="{{ range .spec.containers  }}{{ .name }} {{ end }}"
    __arvan_paas_debug "${FUNCNAME} nouns are ${nouns[@]}"

    local len="${#nouns[@]}"
    if [[ ${len} -ne 1 ]]; then
        return
    fi
    local last=${nouns[${len} -1]}
    local arvan_paas_out
    if arvan_paas_out=$(arvan paas get -o template --template="${template}" pods "${last}" 2>/dev/null); then
        COMPREPLY=( $( compgen -W "${arvan_paas_out[*]}" -- "$cur" ) )
    fi
}

# Require both a pod and a container to be specified
__arvan_paas_require_pod_and_container()
{
    if [[ ${#nouns[@]} -eq 0 ]]; then
        __arvan_paas_parse_get pods
        return 0
    fi;
    __arvan_paas_get_containers
    return 0
}

__custom_func() {
    case ${last_command} in
 
        # first arg is the kind according to ValidArgs, second is resource name
        arvan_paas_get | arvan_paas_describe | arvan_paas_delete | arvan_paas_label | arvan_paas_expose | arvan_paas_export | arvan_paas_patch | arvan_paas_annotate | arvan_paas_edit | arvan_paas_scale | arvan_paas_autoscale | arvan_paas_observe )
            __arvan_paas_get_resource
            return
            ;;

        # first arg is a pod name
        arvan_paas_rsh | arvan_paas_exec | arvan_paas_port-forward | arvan_paas_attach)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __arvan_paas_parse_get pods
            fi;
            return
            ;;
 
        # first arg is a pod name, second is a container name
        arvan_paas_logs)
            __arvan_paas_require_pod_and_container
            return
            ;;
 
        # first arg is a build config name
        arvan_paas_start-build | arvan_paas_cancel-build)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __arvan_paas_parse_get buildconfigs
            fi;
            return
            ;;
 
        # first arg is a deployment config OR deployment
        arvan_paas_rollback)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __arvan_paas_parse_get deploymentconfigs,replicationcontrollers
            fi;
            return
            ;;

        # first arg is a project name
        arvan_paas_project)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __arvan_paas_parse_get projects
            fi;
            return
            ;;
 
        # first arg is an image stream
        arvan_paas_import-image)
            if [[ ${#nouns[@]} -eq 0 ]]; then
                __arvan_paas_parse_get imagestreams
            fi;
            return
            ;;
 
        *)
            ;;
    esac
}
`
)
