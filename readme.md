# PCD

Activities for IF711 - Concurrent and Distributed Programming.

Exercises written in Golang to evaluate the performance of common distributed and concurrent programming concepts such as Mutexes, Channels, TCP, gRPC and RabbitMQ.

The application revolves around a webcrawler that retrieves all url links from a given website. Each folder implements the app but using different programming techniques. 

- `pcd_1/` implements the application concurrently with goroutines.
- `pcd_2/` implements the application sequentially
- `pcd_3/` implements 2 versions of a client-server paradigm: one with tcp and another with udp.
- `pcd_4/` implements 2 verions of a client-server paradigm: one with goRPC and another with gRPC.

`run_100_times.sh` is an utility to help run the applications 100 times and measure the duration of each run. This enables performance comparisons to be made with different implementations.

`analyze_results.ipynb` is a python notebook utility to analyze the data and generate graphs for presentation purposes.


