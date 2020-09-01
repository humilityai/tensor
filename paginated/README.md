# `paginated` tensors

## Overview

Allows for arbitrarily large tensors that primarily live on disk but are buffered in memory. The size of the in-memory buffer can be specified by a page size (in values) and the number of pages allowed in memory. Three default buffer size options are available: small (1/200 of system memory), medium (1/20 of system memory), and large (1/2 of system memory). But any buffer size less than system memory size can be specified.

This package is not ready for production usage or release. Use at your own risk. The API will try to be maintained, but the only API guarantees at this time are the Gorgonia `tensor.Tensor` API.

No behavior is currently guaranteed until a `1.0.0` architecture has been decided on, tested, benchmarked, profiled, and released.

## Contribute

If you would like to contribute then please contribute: tests, benchmarks, documentation, and concurrency-safe code suggestions.

## TODO

- tests, benchmarks, profiling
- concurrency-safe code
