#!/usr/bin/env bats

@test "reject because namespace is not injected" {
  run kwctl run annotated-policy.wasm -r test_data/namespace-disabled.json --settings-json '{"excluded_namespaces": ["bar"]}'

  # this prints the output when one the checks below fails
  echo "output = ${output}"

  # request rejected
  [ "$status" -eq 0 ]
  [ $(expr "$output" : '.*allowed.*false') -ne 0 ]
  [ $(expr "$output" : ".*The 'foo' namespace is not Istio enabled.*") -ne 0 ]
}
