module KV.View where

import KV.Model exposing (Model, KV)
import KV.Update exposing (Action)

import Html exposing (..)
import Html.Attributes exposing (class)
import Signal exposing (Signal, Address)


kvRow : KV -> Html
kvRow kv =
  tr [] [
     td [] [text (toString kv.key)]
    ,td [] [text kv.value]
  ]


view : Signal.Address Action -> Model -> Html
view address model =
  div [class "container-fluid"] [
        h1 [] [text "KVs" ]
      , table [class "table table-striped"] [
          thead [] [
            tr [] [
               th [] [text "key"]
              ,th [] [text "Value"]
          ]
        ]
      , tbody [] (List.map kvRow model.kvs)
    ]
  ]
