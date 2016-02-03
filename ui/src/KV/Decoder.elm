module KV.Decoder where

import Json.Decode as Json exposing ((:=))
import KV.Model exposing (Model, KV)

decode : Json.Decoder (List KV)
decode =
  let
    kv : Json.Decoder KV
    kv =
      Json.object2 KV
        ("key" := Json.string)
        ("value" := Json.string)
  in
    Json.list kv
