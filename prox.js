"use strict";
import fs from "fs"

// 95.181.172.51:8000:tDa005:QXcK8M - ДО
// tDa005:QXcK8M@95.181.172.51:8000 - ПОСЛЕ

const proxy = () => {
  const str = String(fs.readFileSync("./Proxy.txt"))
  const proxs = str.split("\r\n")
  const norm = proxs.map( item => {
   const els =  item.split(":")
   return `${els[2]}:${els[3]}@${els[0]}:${els[1]}`
  } )
  const result = norm.join("\r\n")
  fs.writeFileSync("./Proxy.txt", result)
}
proxy()