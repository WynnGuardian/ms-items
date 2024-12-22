const fetch = require('node-fetch')
const fs = require('fs');
let items = [];

const defaultInt = (some) => some == undefined ? 0 : some;
const defaultDamArray = (some) => some == undefined ? [0, 0] : [some.min, some.max];
let i = 0;
fetch("https://api.wynncraft.com/v3/item/database?fullResult")
.then(x => x.json())
.then(json => {
  var set = new Set();
  Object.keys(json).forEach(key => {
    i++;
    const value = json[key];
    if(!["weapon", "armour", "accessory"].includes(value.type)) {
      return;
    }
    if(value.identifications == undefined) {
      return;
    }
    if(set.has(value.internalName.toLowerCase())) {
      return;
    }
    if(value.icon == undefined) {
      return
    }

    set.add(value.internalName.toLowerCase());
    var item = {
      rarity: value.rarity,
      internalName: key.replace(/[^a-zA-Z0-9\s-]/g, ''),
      levelReq: value.requirements.level,
      strReq: defaultInt(value.requirements.strength),
      agiReq: defaultInt(value.requirements.agility),
      intReq: defaultInt(value.requirements.intelligence),
      dexReq: defaultInt(value.requirements.dexterity),
      defReq: defaultInt(value.requirements.defence),
      sprite: value.icon.value.name,
      baseHP: value.base == undefined ? 0 : defaultInt(value.base.baseHealth),
      baseEarthDef: value.base == undefined ? 0 : defaultInt(value.base.baseEarthDefence), 
      baseAirDef: value.base == undefined ? 0 : defaultInt(value.base.baseAirDefence), 
      baseFireDef:value.base == undefined ? 0 :  defaultInt(value.base.baseFireDefence), 
      baseWaterDef: value.base == undefined ? 0 : defaultInt(value.base.baseWaterDefence), 
      baseThunderDef: value.base == undefined ? 0 : defaultInt(value.base.baseThunderDefence), 
      baseDamage: value.base == undefined ? [0,0] : defaultDamArray(value.base.baseDamage),
      baseEarthDam: value.base == undefined ? [0,0] : defaultDamArray(value.base.baseEarthDamage), 
      baseAirDam: value.base == undefined ? [0,0] : defaultDamArray(value.base.baseAirDamage), 
      baseFireDam: value.base == undefined ? [0,0] : defaultDamArray(value.base.baseFireDamage), 
      baseWaterDam: value.base == undefined ? [0,0] : defaultDamArray(value.base.baseWaterDamage), 
      baseThunderDam: value.base == undefined ? [0,0] : defaultDamArray(value.base.baseThunderDamage), 

      identifications: []
    }
    let allOk = false;
    Object.keys(value.identifications).forEach(id => {
      if(typeof value.identifications[id] === 'object' && value !== null) {
        allOk = true;
      }
      item.identifications = [...item.identifications, {
        stat: id,
        max: value.identifications[id].max,
        min: value.identifications[id].min
      }]
    })
    if(allOk) {
      items = [...items, item]
    }
  });
  fs.writeFile('gen/sanitized.json', JSON.stringify(items), (err) => {
    if (err) {
        console.error('Error writing to file', err);
    } else {
        console.log('JSON data has been written to sanitized.json');
    }
});
})