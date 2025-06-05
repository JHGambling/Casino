/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_719040989")

  // add field
  collection.fields.addAt(2, new Field({
    "hidden": false,
    "id": "bool3978643748",
    "name": "dead",
    "presentable": false,
    "required": false,
    "system": false,
    "type": "bool"
  }))

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_719040989")

  // remove field
  collection.fields.removeById("bool3978643748")

  return app.save(collection)
})
