package main

import (
        "fmt"
        "os"
)

func usage() {
        s := "s3blast - blast files in and out of S3\n"
        s += "USAGE: s3blast [options] OPERATION BUCKET [PREFIX]"
        os.Stderr.WriteString(s)
        os.Exit(0)
}

func Put(bucket string, prefix string) {
        fmt.Println("putting", bucket, prefix)
}

func main() {
        if len(os.Args) < 3 {
                usage()
        }

        op, bucket := os.Args[1], os.Args[2]
        prefix := ""
        if len(os.Args) >= 4 {
                prefix = os.Args[3]
        }

        switch op {
        case "LIST":
                List(bucket, prefix)
        case "PUT":
                Put(bucket, prefix)
        default:
                usage()
        }
}
