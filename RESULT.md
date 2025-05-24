

```PROJECT 1 Result

teotia@mac tinykv % make project1
GO111MODULE=on go test -v --count=1 --parallel=1 -p=1 ./kv/server -run 1
=== RUN   TestRawGet1
--- PASS: TestRawGet1 (1.00s)
=== RUN   TestRawGetNotFound1
--- PASS: TestRawGetNotFound1 (0.58s)
=== RUN   TestRawPut1
--- PASS: TestRawPut1 (0.90s)
=== RUN   TestRawGetAfterRawPut1
--- PASS: TestRawGetAfterRawPut1 (0.70s)
=== RUN   TestRawGetAfterRawDelete1
--- PASS: TestRawGetAfterRawDelete1 (0.76s)
=== RUN   TestRawDelete1
--- PASS: TestRawDelete1 (0.38s)
=== RUN   TestRawScan1
--- PASS: TestRawScan1 (0.75s)
=== RUN   TestRawScanAfterRawPut1
--- PASS: TestRawScanAfterRawPut1 (1.09s)
=== RUN   TestRawScanAfterRawDelete1
--- PASS: TestRawScanAfterRawDelete1 (0.67s)
=== RUN   TestIterWithRawDelete1
--- PASS: TestIterWithRawDelete1 (1.09s)
PASS
ok  	github.com/pingcap-incubator/tinykv/kv/server	10.447s

```