Specification
---

```yaml
version:
  0:
    transports:
      - type: http
        versions: ["1.0", "1.1", "2"]
        mappings:
          - prefix: "/api/v{{$version}}"
            mappings:
              - path: "/buckets"
                call:
                  name: GetBuckets
              - path: "/buckets/{id}"
                call:
                  name: GetBucket
                  arguments:
                    - name: id
                      value: "$path.id"
    functions:
      - name: GetBuckets
        returns:
          - name: buckets
            type: "[]Bucket"
      - name: GetBucket
        arguments:
          - name: id
            type: string
        returns:
          - name: bucket
            type: Bucket
    types:
      - name: Bucket
        fields:
          - name: id
            type: string
          - name: name
            type: string
```
