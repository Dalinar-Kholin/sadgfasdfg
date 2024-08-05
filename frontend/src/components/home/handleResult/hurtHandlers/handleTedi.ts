import {IHurtInfoForComp} from "../handleResultInterfaces.ts";
import {hurtNames} from "../../../../interfaces.ts";
import round from "./handleSot.ts";

export function handleTedi() {
    return (serverResponse: any) : IHurtInfoForComp => {
        const noProd = {
            hurtName: hurtNames.tedi,
            priceForPack: -1,
            priceForOne: -1,
            productsInPack: -1
        }
        if (serverResponse.response === null
            || serverResponse.count === 0
            || (serverResponse.results === null || serverResponse.results === undefined || serverResponse.results.length === 0)) {
            return noProd
        }else{


            const res = serverResponse.results.reduce((min : any,curr: any) => {
                return (+(curr.final_price) < +(min.final_price)) && (curr.stocks[0].quantity_available > 0) ? curr : min
            })


            if (res.stocks[0].quantity_available === 0) {
                return noProd
            }
            const pack = +(res.cumulative_unit_ratio_splitter)
            const priceForOneItem = +(res.final_price);
            const productInPack = (Math.ceil(1 / pack) * pack);
            return {
                hurtName: hurtNames.tedi,
                priceForOne: priceForOneItem,
                productsInPack: productInPack,
                priceForPack: round(priceForOneItem * productInPack)
            }
        }
    }
}
