console.log("Script received");

function koimetricsIsPhoneDevice() {
    if (navigator.userAgent.match(/Android/i)
        || navigator.userAgent.match(/webOS/i)
        || navigator.userAgent.match(/iPhone/i)
        || navigator.userAgent.match(/iPad/i)
        || navigator.userAgent.match(/iPod/i)
        || navigator.userAgent.match(/BlackBerry/i)
        || navigator.userAgent.match(/Windows Phone/i)) {
        return 1;
    } else {
        return 0;
    }
}

function koimetricsFormatDate(date) {
    var d = new Date(date),
    month = '' + (d.getMonth() + 1);
    day = '' + d.getDate();
    year = d.getFullYear();    
    month = (month.length < 2)? '0'+month:month;
    day = (day.length < 2 )? '0'+day:day;
    return [year, month, day].join('-');
}


function koimetricsFormatHHMM(date){
    var d = new Date(date);
    hh = '' + d.getHours(),
    hh = (hh.length < 2)? '0'+hh:hh;
    mm = '' + d.getMinutes();
    mm = (mm.length < 2)? '0'+mm:mm;
    ss = '' + d.getSeconds();
    ss = (ss.length < 2)? '0'+ss:ss;
    return [hh,mm,ss].join(':');
}

function heart_beat( ){
    var f_data = new FormData();
    f_data.append("Key", "{{.key}}");
    f_data.append("session_id", "{{.session_id}}");
    fetch("{{.goapi_host}}/API/v1/heartbeats/", {
        method: "POST",
        body: f_data,
    }).then(function (res) {
        console.log(res.json);
    });
}

function koimetricsSendData(koimetricsData) {
    var f_data = new FormData();
    for (var key in koimetricsData) {
        f_data.append(key, koimetricsData[key]);
    }
    fetch("{{.goapi_host}}/API/v1/statistics/", {
        method: "POST",
        body: f_data,
    }).then(function (res) {
        console.log(res.json);
    });
    console.log("Data sent");
    window.setInterval(heart_beat, 10000);
}

function captureUserStatus() {
    var data = new Object();
    var date = new Date();
    data["Key"]  = "{{.key}}";
    data["Host"] = window.location.host;
    data["Path"] = window.location.pathname;
    data["Date"] = koimetricsFormatDate(date);
    data["Referrer"] = (document.referrer.length > 0 ? document.referrer.split("/")[2].replace("www.", "") : ""),
    data["ReferrerPath"] = (data["Referrer"].length > 0 ? "/"+document.referrer.split("/").slice(3).join("/") : ""),
    data["Time"] = koimetricsFormatHHMM(date)
    data["Performance"] = performance.now();
    data["Latitude"] = "";
    data["Longitude"] = "";
    data["IsPhone"] = koimetricsIsPhoneDevice();
    data["session_id"] = "{{.session_id}}";

    // Fetch get users location data from IPAPI
    fetch('https://ipapi.co/json')
        .then(response => response.json())
        .then(jresponse => {
            var country = jresponse.country_name;
            var city = jresponse.city;
            var regionName = jresponse.region;
            var latitude = jresponse.latitude;
            var longitude = jresponse.longitude;
            data["Country"] = country;
            data["City"] = city;
            data["Region"] = regionName;
            data["Latitude"] = latitude;
            data["Longitude"] = longitude;
            
            var askLocationTo = "{{.ask_location_to}}";
            if (askLocationTo.includes(window.location.host)) {
                if (navigator.geolocation) {
                    navigator.geolocation.getCurrentPosition((position) => {
                        data["Latitude"] = position.coords.latitude;
                        data["Longitude"] = position.coords.longitude;
                    });
                }
            }
        })
        .catch(function (err) {
            var askLocationTo = "{{.ask_location_to}}";
            if (askLocationTo.includes(window.location.host)) {
                if (navigator.geolocation) {
                    navigator.geolocation.getCurrentPosition((position) => {
                        data["Latitude"] = position.coords.latitude;
                        data["Longitude"] = position.coords.longitude;
                    });
                    koimetricsSendData(data);
                }
               // Send data without user location
            }
        });
    // Check users location
    setTimeout(function(){
        koimetricsSendData(data);
    }, 3000);
}
captureUserStatus();