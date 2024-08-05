import {IHurtInfoForComp} from "../handleResultInterfaces.ts";
import {hurtNames} from "../../../../interfaces.ts";
import round from "./handleSot.ts";




export function handleSpecialAndSot(name: hurtNames) {
    return (serverResponse: any) : IHurtInfoForComp => {
        if (serverResponse === null || serverResponse=== -1 ){
            return {
                hurtName: name,
                priceForPack: -1,
                priceForOne: -1,
                productsInPack: -1
            }
        }else{
            const pack = serverResponse.pozycje[0].ilOpkZb;
            const priceForOneItem = serverResponse.pozycje[0].cenaNettoOstateczna;
            const productInPack = (Math.ceil(1 / pack) * pack);
            return {
                hurtName: name,
                priceForOne: priceForOneItem,
                productsInPack: productInPack,
                priceForPack: round(priceForOneItem * productInPack)
            }
        }
    }
}