FROM prod-registry.cicd.prod.core.n-t.io/repository/docker/core/library/toolbox/toolbox-test-go:master as builder
COPY ./ /src
RUN \
    for dir in $(find /src/test/chaincode/public -mindepth 2 -maxdepth 2 -type d); \
    do \
        cd "$dir" && echo "$dir" && go mod tidy && go mod vendor ; \
    done

FROM prod-registry.cicd.prod.core.n-t.io/repository/docker/core/library/toolbox/toolbox-test-go:master
COPY --from=builder /src/test/chaincode /chaincode
