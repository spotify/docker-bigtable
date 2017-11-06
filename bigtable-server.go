/*
 * Copyright 2016 Spotify AB.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */
package main

import (
        "context"
        "fmt"
        "os"
        "os/signal"
        "syscall"
        "cloud.google.com/go/bigtable/bttest"
        "cloud.google.com/go/bigtable"
        "google.golang.org/grpc"
        "strings"
        "google.golang.org/api/option"
        "flag"
        "errors"
)

func main() {
        cfs := flag.String("cf", "", "Optional: the column families to create at startup. Format: a series of <instance>.<table>.<column familiy>, comma separated. Ex: \n          docker run -d spotify/bigtable-emulator -cf dev.records.data,dev.records.metadata")
        flag.Parse()

        srv, err := bttest.NewServer("0.0.0.0:8080")
        if err != nil {
                fmt.Fprintf(os.Stderr, "Error starting the server: %v\n", err)
                return
        }
        defer srv.Close()

        if err = createColumnFamiliies(*cfs); err != nil {
                fmt.Fprintf(os.Stderr, "Error creating the column familiies: %v\n", err)
                return
        }

        sigs := make(chan os.Signal, 1)
        signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

        fmt.Printf("bttest.Server running on %s\n", srv.Addr)
        <-sigs
        fmt.Println("done")
}

func createColumnFamiliies(specifications string) error {
        if specifications == "" {
                return nil
        }

        ctx := context.Background()
        conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
        if err != nil {
                return fmt.Errorf("could not connect to bigtable emulator: %v", err)
        }

        for _, specification := range strings.Split(specifications, ",") {
                specificationElements := strings.Split(specification, ".")
                if len(specificationElements) != 3 {
                        return errors.New("format of column family to create is <instance>.<table>.<column family>")
                }

                instance := specificationElements[0]
                table := specificationElements[1]
                columnFamily := specificationElements[2]

                client, err := bigtable.NewAdminClient(ctx, "dev", instance, option.WithGRPCConn(conn))
                if err != nil {
                        return fmt.Errorf("failed to create admin client: %v", err)
                }

                tables, err := client.Tables(ctx)
                if !tableExists(tables, table) {
                        if err = client.CreateTable(ctx, table); err != nil {
                                return err
                        }
                }

                tableInfo, err := client.TableInfo(ctx, table)
                if !columnFamilyExists(tableInfo.FamilyInfos, columnFamily) {
                        fmt.Printf("creating %v.%v.%v column family\n", instance, table, columnFamily)
                        if err := client.CreateColumnFamily(ctx, table, columnFamily); err != nil {
                                return err
                        }
                }
        }

        return nil
}

func tableExists(tables []string, table string) bool {
        for _, item := range tables {
                if item == table {
                        return true
                }
        }
        return false
}

func columnFamilyExists(columnFamilies []bigtable.FamilyInfo, columnFamily string) bool {
        for _, family := range columnFamilies {
                if family.Name == columnFamily {
                        return true
                }
        }
        return false
}
