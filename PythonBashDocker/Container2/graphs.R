library(lattice)

read.table("output-throughput-latency/stats.csv", header=TRUE) -> csvDataFrameSource
csvDataFrame <- csvDataFrameSource

trellis.device("pdf", file="graph1.pdf", color=T, width=6.5, height=5.0)

xyplot(requests ~ rate, 
       data = csvDataFrame,
       type = "b",
       xlab = "Rate",
       ylab = "Response/Throughput")


dev.off() -> null 

trellis.device("pdf", file="graph2.pdf", color=T, width=6.5, height=5.0)
xyplot(latency ~ requests,
       data = csvDataFrame,
       type = "b",
       ylab = "Latency",
       xlab = "Response/Throughput")

dev.off() -> null 
