(function () {
    "use strict";
    const forms = document.querySelectorAll(".requires-validation");
    Array.from(forms).forEach(function (form) {
        form.addEventListener(
            "submit",
            function (event) {
                if (!form.checkValidity()) {
                    event.preventDefault();
                    event.stopPropagation();
                }

                form.classList.add("was-validated");
            },
            false
        );
    });
})();

function functionToAsk() {
    emptydiv = "<span> </span>";
    document.getElementById("llmanswer").innerHTML = emptydiv;

    dd = document.getElementById("query");

    var myLLM = $('input[name="LLM"]:checked').val();

    var taskArray = new Array();
    $("input[name='LLM']").each(function() {
        if (this.checked) {
            taskArray.push($(this).val());
        }
    });

    myLLM = taskArray;

    var preprompt = "";
    var postprompt = "";

    if ($('#prompts').is(':checked')) {
        preprompt = $("#preprompt").val();
        postprompt = $("#postprompt").val();
    }

    var request = {};
    request.query = dd?.value;
    request.engines = myLLM;
    request.pre = preprompt;
    request.post = postprompt;

    jsondata = JSON.stringify(request);

    if (request.engines[0] === undefined) {
        $("#engine").show(true);
        $("#inValidFeedback").show(true);
    } else {
        $("#engine").show(true);
        $("#inValidFeedback").hide(true);
        $("#validFeedback").show(true);
    }

    // alert(jsondata);
    // $("#overlay").show();

    var url = "http://localhost:8080/query";
    var url2 =
        "https://ttec-vk-site-dot-insightsteamsandbox.uc.r.appspot.com/query";

    var success = false;
    var Response = {};

    $.ajax({
        url: url2,
        crossDomain: true,
        method: "POST",
        contentType: "application/json",
        data: jsondata,
        headers: { "Access-Control-Allow-Origin": "*" },
        dataType: "json",
        type: "POST",
        async: true,

        beforeSend: function () {
            // Show the loader before the request is sent
            $("#overlay").show();
        },

        success: function (Data, textStatus, jqXHR) {
            success = true;
            Response = Data;
            SuccessCallback(Response);
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            alert("Status: " + textStatus);
            alert("Error: " + errorThrown);
            $("#overlay").fadeOut();
        },

        // complete: function () {
        //   // Hide the loader after the request is complete
        //   $("#overlay").hide();
        // },
    });
}

function SuccessCallback(Response) {

    $("#answer").show(true);
    $("#overlay").fadeOut();

    var newdiv = "";
    Response.results.forEach((element) => {
        var citation = '  <div class="row m-2 mt-4 mr-4" id="answer" style="display: block">'
        citation = citation + '<div class="col-md-12">'
        citation = citation + '<div class="card shadow p-3 bg-body rounded">'
        citation = citation + '    <div class="card-body style="overflow-wrap: break-word">'
        
        // citation = citation + ' <button type="button" class="btn btn-outline-dark float-end citationButton" onClick="showCitation()">Citation</button>';

        var nl = '<span style="white-space: pre-wrap;" data-id="';
        nl = nl + element.llmname + '">' + element.llmname + " / " + element.llmcompute + "ms </p>" +
            element.llmresponse;
        nl = nl + "</p></span>";

        nl = nl + "<p> </p><b> </b>"
        if (element.llmcitations.length != null) {
            element.llmcitations.forEach((cit) => {

                nl = nl + '<a href="https://storage.cloud.google.com/llm-serve/viking/' + cit + '" target="_blank" style="font-size: 9px; color: blue;">' + cit + '</a>';
                nl = nl + '<p style="margin: 0;"> </p>';
            });
        }

        newdiv = newdiv + citation +  nl;
        newdiv = newdiv +'    </div></div></div></div>'

    });
    document.getElementById("llmanswer").innerHTML = newdiv;


}

var citationDiv = "";

function showCitation() {
    $("#citationDiv").show(true);
    var citation = "this is citation div";
    citationDiv = citation + citationDiv;
    document.getElementById("citation").innerHTML = citationDiv;
}

function CustomPromptsValueChaned(element) {
    if($('#prompts').is(':checked')) {
        $("#llmPrompts").show(true);
    } else {
        $("#llmPrompts").hide();
    }
}

// $("#overlay")
//   .hide() // Hide it initially
//   .ajaxStart(function () {
//     $(this).show();
//   })
//   .ajaxStop(function () {
//     $(this).hide();
//   });

$('input[name="LLM"]').on("change", function () {
    var myLLM = $('input[name="LLM"]:checked').val();
    if (myLLM !== null || myLLM !== undefined) {
        $("#engine").show(true);
        $("#inValidFeedback").hide(true);
        $("#validFeedback").show(true);
    } else {
        $("#engine").show(true);
        $("#inValidFeedback").show(true);
        $("#validFeedback").hide(true);
    }
});