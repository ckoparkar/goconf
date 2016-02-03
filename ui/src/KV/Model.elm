module KV.Model where

import String
import Dict


type alias KV =
  { key : String
  , value : String
  }


fromTuple (x,y) =
  KV x y


folder : KV -> Bool
folder kv =
  String.contains "/" kv.key


type alias Model =
  { kvs : List KV
  }


-- testKVs = Model [(KV "hello/a" "world"), (KV "hello/b" "world"), (KV "abc" "abc")]

uniqueKVs : List KV -> List KV
uniqueKVs kvs =
  let
    splitKey : KV -> (String, String)
    splitKey kv =
      ( (String.split "/" kv.key
        |> List.head
        |> Maybe.withDefault "")
      , kv.value)
  in
    List.map (\kv -> if folder kv
                     then (KV kv.key "")
                     else kv) kvs
      |> List.map splitKey
      |> Dict.fromList
      |> Dict.toList
      |> List.map fromTuple
