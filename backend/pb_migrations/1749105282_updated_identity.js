/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_719040989")

  // add field
  collection.fields.addAt(1, new Field({
    "cascadeDelete": true,
    "collectionId": "pbc_1266413393",
    "hidden": false,
    "id": "relation2087227935",
    "maxSelect": 1,
    "minSelect": 0,
    "name": "wallet",
    "presentable": false,
    "required": true,
    "system": false,
    "type": "relation"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_719040989")

  // remove field
  collection.fields.removeById("relation2087227935")

  return app.save(collection)
})
