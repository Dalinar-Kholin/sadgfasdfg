import {IItemToSearch} from "../../../../interfaces.ts";
import handlePcMarket from "./PcMarket.ts";
import handleKc from "./handleKc.ts";

export default function stratchInputType(i :string, expansion: string){ // sprawdza jakiego rodzaju jest plik wejściowy a następnie wybiera dla niego odpowiednią funkcję
    const lmbd = () : IItemToSearch[] => {
        return [{Ean: "0", Amount: 0, Name: ""}]}
    switch (expansion.substring(expansion.lastIndexOf('.') + 1, expansion.length).toLowerCase()){
        case "txt":
            if (i.includes("Linia:Nazwa")){
                return handlePcMarket
            }else{
                return handleKc
            }
        case "pdf":
            return lmbd
        case "xlsx":
            return  lmbd
    }
    return lmbd
}
