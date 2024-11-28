<?php
    $TEST_CASES = '[{"input": [1,2], "expected": 3}]';
    function sum(int $a, int $b):int{
        return $a+$b;
    }

    $result = [
        "final" => true,
        "detail" => []
    ];

    try{
        foreach(json_decode($TEST_CASES) as $case){
           $tmp_result = call_user_func("sum", ...$case->input);
           if($case->expected != $tmp_result)
                $result['final'] = false;
           array_push($result['detail'] , [
            "input" => $case->input,
            "expected" => $case->expected,
            "passed" => $case->expected == $tmp_result,
           ]);
        }
    }catch(Exception $error){
        $result['final'] = false;
        array_push($result['detail'] , [
            "input" => $case->input,
            "expected" => $case->expected,
            "passed" => false,
            "err" => $error->getMessage()
        ]);
    }

    echo json_encode($result);
?>