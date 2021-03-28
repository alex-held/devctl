if [[ ! -o interactive ]]; then
    return
fi

compctl -K _devctlenv devctlenv

_devctlenv() {
  local words completions
  read -cA words

  if [ "${#words}" -eq 2 ];hen
    completions="$(devctlenv commands)"
  else
    completions="$(devctlenv completions ${words[2,-2]})"
  fi

  reply=("${(ps:\n:)completions}")
}
