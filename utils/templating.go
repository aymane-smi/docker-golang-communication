package utils

import (
	"aymane/types"
	"encoding/json"
	"strings"
)

func InitTemplateJS(code string, fname string, cases []types.Cases) string {
	template := `
		{code}
		const TEST_CASES = {test_cases}
		let result = {
			final : false,
			detail: []
		}
		for(let case_ of TEST_CASES){
			const {Input, Expected} = case_;
			try{
				const tmp_result = {fname}(...Input)
				result.detail.push({Input, Expected, passed: tmp_result === Expected});
				if(tmp_result !== Expected)
					result.final = false;
			}catch(err){
				result.final = false;
				result.detail.push({Input, Expected, passed: false, error: err.message})
			}
		}
		console.log(JSON.stringify(result));
	`
	finalCode := strings.Replace(template, "{code}", code, -1)
	json, err := json.Marshal(cases)
	if err != nil {
		panic(err.Error())
	}
	finalCode = strings.Replace(finalCode, "{test_cases}", string(json), -1)
	finalCode = strings.Replace(finalCode, "{fname}", fname, -1)
	return finalCode
}
