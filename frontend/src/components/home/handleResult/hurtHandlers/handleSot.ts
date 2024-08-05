/*
export function handleSot() {
    return (serverResponse: any) : IHurtInfoForComp => {
        if (serverResponse.countPozycji === 0){
            return {
                priceForPack: -1,
                priceForOne: -1,
                productsInPack: -1
            }
        }else{
            const pack = serverResponse.pozycje[0].ilOpkZb;
            const priceForOneItem = serverResponse.pozycje[0].cenaNettoOstateczna;
            const productInPack = (Math.ceil(1 / pack) * pack);
            return {
                priceForOne: priceForOneItem,
                productsInPack: productInPack,
                priceForPack: priceForOneItem * productInPack
            }
        }
    }
}*/
export default function round(x: number) {
    return Math.round(x * 100) / 100
}