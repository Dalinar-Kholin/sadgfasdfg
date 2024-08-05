import {hurtNames} from '../../../interfaces.ts';
import {IHurtInfoForComp} from "./handleResultInterfaces.ts";
import {handleEurocash} from "./hurtHandlers/handleEurocash.ts";
import {handleTedi} from "./hurtHandlers/handleTedi.ts";
import {handleSpecialAndSot} from "./hurtHandlers/handleSpecial.ts";


interface IHandleResults {
    name : hurtNames
}


export function handleResults({name }: IHandleResults): (x : any) => IHurtInfoForComp{
    switch (name) {
        case hurtNames.eurocash:
            return handleEurocash();
        case hurtNames.special:
            return handleSpecialAndSot(hurtNames.special);
        case hurtNames.sot:
            return handleSpecialAndSot(hurtNames.sot);
        case hurtNames.tedi:
            return handleTedi();
        default:
            return (_x : any): IHurtInfoForComp=>{
                return {
                    hurtName: hurtNames.none,
                    priceForPack: -1,
                    priceForOne : -1,
                    productsInPack: -1
                }
            };
    }
}

