/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("_pb_users_auth_")

  // remove field
  collection.fields.removeById("relation1788209604")

  // add field
  collection.fields.addAt(9, new Field({
    "cascadeDelete": false,
    "collectionId": "pbc_1266413393",
    "hidden": false,
    "id": "relation2087227935",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "wallet",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "relation"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("_pb_users_auth_")

  // add field
  collection.fields.addAt(8, new Field({
    "cascadeDelete": false,
    "collectionId": "pbc_719040989",
    "hidden": false,
    "id": "relation1788209604",
    "maxSelect": 9000,
    "minSelect": 0,
    "name": "identities",
    "presentable": false,
    "required": true,
    "system": false,
    "type": "relation"
  }))

  // remove field
  collection.fields.removeById("relation2087227935")

  return app.save(collection)
})
