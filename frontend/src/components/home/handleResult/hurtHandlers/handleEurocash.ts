import {IHurtInfoForComp} from "../handleResultInterfaces.ts";
import {hurtNames} from "../../../../interfaces.ts";
import round from "./handleSot.ts";

export function handleEurocash(){
    return (serverResponse: any) : IHurtInfoForComp => {

        if (serverResponse === null ||  serverResponse.Data.TotalCount ===0){
            return {
                hurtName: hurtNames.eurocash,
                priceForPack: -1,
                priceForOne : -1,
                productsInPack: -1
            }
        }else{
            const pack = serverResponse.Data.Items[0].SposobPakowania;
            const priceForOneItem = serverResponse.Data.Items[0].CenaBudzet;
            const productInPack = (Math.ceil(1 / pack) * pack);
            return {
                hurtName: hurtNames.eurocash,
                priceForOne: priceForOneItem,
                productsInPack: productInPack,
                priceForPack: round(priceForOneItem * productInPack * 100) / 100
            }
        }
    }
}

