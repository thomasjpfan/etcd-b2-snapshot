.PHONY: unittest image

TAG?=master
TRAVIS_COMMIT?=master

unittest:
	go test ./... --run UnitTest

image:
	docker image build -t thomasjpfan/etcd-b2-snapshot:$(TAG) \
	--label "org.opencontainers.image.revision=$(TRAVIS_COMMIT)" .
