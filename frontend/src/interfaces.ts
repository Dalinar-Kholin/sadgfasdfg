import {IHurtInfoForComp} from "./components/home/handleResult/handleResultInterfaces.ts";
import {IServerMultipleDataResult} from "./components/home/resultGrabbers.ts";

export type PATH= "main" | "cennik" | "ustawienia" | "login"


export enum hurtNames{
    "none" = 0,
    "eurocash"= 1,
    "special"= 2,
    "sot"= 4,
    "tedi"= 8,
}


export interface IAllResult {
    ean: string,
    result: IServerMultipleDataResult[]
}

export interface IItemToSearch {
    Name : string,
    Ean : string,
    Amount: number
}

export interface IItemInstance {
    item : IHurtInfoForComp
    name : string // brane z pliku tekstowego
    ean: string,
    count: number,
}
