
function sum(a,b){return a+b}
const TEST_CASES = [{"Input":[1,2],"Expected":3}]
let result = {
    final : false,
    detail: []
}
for(let case_ of TEST_CASES){
    const {Input, Expected} = case_;
    try{
        const tmp_result = sum(...Input)
        result.detail.push({Input, Expected, passed: tmp_result === Expected});
        console.log(tmp_result)
        if(tmp_result !== Expected)
            result.final = false;
    }catch(err){
        result.final = false;
        result.detail.push({Input, Expected, passed: false, error: err.message})
    }
}
console.log(JSON.stringify(result));
