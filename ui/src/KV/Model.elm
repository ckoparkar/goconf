module KV.Model where

type alias KV =
  { key : String
  , value : String
  }


type alias Model =
  { kvs : List KV
  }
