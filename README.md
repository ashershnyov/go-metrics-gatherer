# go-musthave-metrics-tpl

## Инкремент 17: разница в аллокациях после оптимизации сжимающей middleware и чтения JSON в хендлерах:
Для нагрузки использовалась утилита [oha](https://github.com/hatoo/oha)

### Без ограничения RPS:

До оптимизации:
```bash
Summary:
  Success rate: 100.00%
  Total:        40002.5820 ms
  Slowest:      276.6544 ms
  Fastest:      0.0615 ms
  Average:      3.6683 ms
  Requests/sec: 13619.4959

  Total data:   11.95 MiB
  Size/request: 23 B
  Size/sec:     305.89 KiB

Response time histogram:
    0.061 ms [1]      |
   27.721 ms [542676] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
   55.380 ms [1324]   |
   83.039 ms [492]    |
  110.699 ms [187]    |
  138.358 ms [47]     |
  166.017 ms [4]      |
  193.676 ms [0]      |
  221.336 ms [0]      |
  248.995 ms [0]      |
  276.654 ms [50]     |

Response time distribution:
  10.00% in 0.2565 ms
  25.00% in 0.5495 ms
  50.00% in 1.5406 ms
  75.00% in 5.3346 ms
  90.00% in 8.8704 ms
  95.00% in 11.3132 ms
  99.00% in 18.6355 ms
  99.90% in 68.4807 ms
  99.99% in 138.1481 ms


Details (average, fastest, slowest):
  DNS+dialup:   1.6837 ms, 1.1934 ms, 2.5302 ms
  DNS-lookup:   0.0502 ms, 0.0007 ms, 0.4489 ms

Status code distribution:
  [200] 544781 responses
```

После оптимизации:
```bash
Summary:
  Success rate: 100.00%
  Total:        40003.0780 ms
  Slowest:      280.9724 ms
  Fastest:      0.0279 ms
  Average:      1.3379 ms
  Requests/sec: 37230.0851

  Total data:   32.67 MiB
  Size/request: 23 B
  Size/sec:     836.22 KiB

Response time histogram:
    0.028 ms [1]       |
   28.122 ms [1486738] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
   56.217 ms [1530]    |
   84.311 ms [541]     |
  112.406 ms [201]     |
  140.500 ms [104]     |
  168.595 ms [1]       |
  196.689 ms [36]      |
  224.784 ms [64]      |
  252.878 ms [0]       |
  280.972 ms [100]     |

Response time distribution:
  10.00% in 0.1675 ms
  25.00% in 0.3608 ms
  50.00% in 0.9095 ms
  75.00% in 1.5995 ms
  90.00% in 2.4050 ms
  95.00% in 3.1352 ms
  99.00% in 7.5136 ms
  99.90% in 46.0616 ms
  99.99% in 203.2781 ms


Details (average, fastest, slowest):
  DNS+dialup:   2.1520 ms, 1.8202 ms, 2.9511 ms
  DNS-lookup:   0.0879 ms, 0.0007 ms, 0.6694 ms

Status code distribution:
  [200] 1489316 responses
```

Разница отчетов pprof:
```bash
File: main
Type: alloc_space
Time: 2026-02-04 20:20:01 MSK
Showing nodes accounting for -607.06GB, 98.93% of 613.63GB total
Dropped 192 nodes (cum <= 3.07GB)
      flat  flat%   sum%        cum   cum%
 -336.57GB 54.85% 54.85%  -609.19GB 99.28%  compress/flate.NewWriter (inline)
 -166.41GB 27.12% 81.97%  -272.62GB 44.43%  compress/flate.(*compressor).init
 -104.20GB 16.98% 98.95%  -104.20GB 16.98%  compress/flate.newDeflateFast (inline)
    0.09GB 0.015% 98.93%  -607.01GB 98.92%  github.com/ashershnyov/go-metrics-gatherer/internal/server.New.Logging.func1.1
    0.01GB 0.0022% 98.93%  -607.07GB 98.93%  github.com/ashershnyov/go-metrics-gatherer/internal/server.New.Gzip.func2.1
    0.01GB 0.0018% 98.93%  -606.36GB 98.82%  net/http.(*conn).serve
         0     0% 98.93%  -609.19GB 99.28%  compress/gzip.(*Writer).Close
         0     0% 98.93%  -609.19GB 99.28%  compress/gzip.(*Writer).Write
         0     0% 98.93%  -606.78GB 98.88%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 98.93%  -607.01GB 98.92%  net/http.HandlerFunc.ServeHTTP
         0     0% 98.93%  -606.78GB 98.88%  net/http.serverHandler.ServeHTTP
```

### С ограничением 2000 RPS:

До оптимизации:
```bash
Summary:
  Success rate: 100.00%
  Total:        40002.6972 ms
  Slowest:      329.5020 ms
  Fastest:      0.0918 ms
  Average:      2.3284 ms
  Requests/sec: 1999.8401

  Total data:   1.75 MiB
  Size/request: 23 B
  Size/sec:     44.92 KiB

Response time histogram:
    0.092 ms [1]     |
   33.033 ms [78568] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
   65.974 ms [1074]  |
   98.915 ms [219]   |
  131.856 ms [77]    |
  164.797 ms [6]     |
  197.738 ms [10]    |
  230.679 ms [20]    |
  263.620 ms [14]    |
  296.561 ms [9]     |
  329.502 ms [1]     |

Response time distribution:
  10.00% in 0.2759 ms
  25.00% in 0.3913 ms
  50.00% in 0.5493 ms
  75.00% in 0.8090 ms
  90.00% in 1.4502 ms
  95.00% in 7.8427 ms
  99.00% in 45.2100 ms
  99.90% in 122.4378 ms
  99.99% in 273.1375 ms


Details (average, fastest, slowest):
  DNS+dialup:   0.3805 ms, 0.0757 ms, 4.6281 ms
  DNS-lookup:   0.0147 ms, 0.0008 ms, 0.3103 ms

Status code distribution:
  [200] 79999 responses
```

После оптимизации:
```bash
Summary:
  Success rate: 100.00%
  Total:        40004.1399 ms
  Slowest:      57.1388 ms
  Fastest:      0.0533 ms
  Average:      0.7209 ms
  Requests/sec: 1999.7680

  Total data:   1.75 MiB
  Size/request: 23 B
  Size/sec:     44.92 KiB

Response time histogram:
   0.053 ms [1]     |
   5.762 ms [78062] |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  11.470 ms [677]   |
  17.179 ms [330]   |
  22.887 ms [328]   |
  28.596 ms [296]   |
  34.305 ms [129]   |
  40.013 ms [127]   |
  45.722 ms [38]    |
  51.430 ms [3]     |
  57.139 ms [7]     |

Response time distribution:
  10.00% in 0.1294 ms
  25.00% in 0.1663 ms
  50.00% in 0.2170 ms
  75.00% in 0.2857 ms
  90.00% in 0.4407 ms
  95.00% in 0.8786 ms
  99.00% in 18.7710 ms
  99.90% in 39.3777 ms
  99.99% in 48.5976 ms


Details (average, fastest, slowest):
  DNS+dialup:   0.2522 ms, 0.0785 ms, 1.2416 ms
  DNS-lookup:   0.0417 ms, 0.0009 ms, 0.4833 ms

Status code distribution:
  [200] 79998 responses
```

Разница отчетов pprof:
```bash

File: main
Type: alloc_space
Time: 2026-02-04 20:28:04 MSK
Showing nodes accounting for -91179.91MB, 97.89% of 93146.94MB total
Dropped 168 nodes (cum <= 465.73MB)
      flat  flat%   sum%        cum   cum%
-50657.64MB 54.38% 54.38% -91500.15MB 98.23%  compress/flate.NewWriter (inline)
-24929.88MB 26.76% 81.15% -40842.51MB 43.85%  compress/flate.(*compressor).init
-15594.38MB 16.74% 97.89% -15594.38MB 16.74%  compress/flate.newDeflateFast (inline)
    2.50MB 0.0027% 97.89% -91916.49MB 98.68%  github.com/ashershnyov/go-metrics-gatherer/internal/server.New.Logging.func1.1
   -0.50MB 0.00054% 97.89% -91590.70MB 98.33%  github.com/ashershnyov/go-metrics-gatherer/internal/server.New.Gzip.func2.1
         0     0% 97.89%  -463.09MB   0.5%  bufio.(*Writer).Flush
         0     0% 97.89% -91500.15MB 98.23%  compress/gzip.(*Writer).Close
         0     0% 97.89% -91500.15MB 98.23%  compress/gzip.(*Writer).Write
         0     0% 97.89% -91966.03MB 98.73%  github.com/go-chi/chi/v5.(*Mux).ServeHTTP
         0     0% 97.89%  -463.09MB   0.5%  net/http.(*chunkWriter).Write
         0     0% 97.89%  -463.09MB   0.5%  net/http.(*chunkWriter).writeHeader
         0     0% 97.89% -92577.85MB 99.39%  net/http.(*conn).serve
         0     0% 97.89%  -502.13MB  0.54%  net/http.(*response).finishRequest
         0     0% 97.89% -91916.49MB 98.68%  net/http.HandlerFunc.ServeHTTP
         0     0% 97.89% -91966.03MB 98.73%  net/http.serverHandler.ServeHTTP
         0     0% 97.89%  -845.45MB  0.91%  sync.(*Pool).Get
```