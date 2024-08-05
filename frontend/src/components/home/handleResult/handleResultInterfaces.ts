import {hurtNames} from "../../../interfaces.ts";

export interface IHurtInfoForComp {
    hurtName: hurtNames
    priceForPack: number,
    priceForOne: number,
    productsInPack: number
}
