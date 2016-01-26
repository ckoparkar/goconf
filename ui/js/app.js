var url = "http://localhost:8080";

$(document).on({
    ajaxStart: function() { $('#ajaxloader').fadeIn();    },
    ajaxStop: function() { $('#ajaxloader').fadeOut(); }
});

function getAllKVs(token) {
    $('#ajaxloader').fadeIn();
    var kvs = [];
    $.ajax( url + "/v1/kv/?recurse=true", {
        type: "get",
        async: false,
        data: {token: token},
        success: function(raw) {
            kvs = JSON.parse(raw);
            kvs = $.map(kvs, function(item) { return new KV(item.key, item.value); });
            $('#ajaxloader').fadeOut();
        }
    });
    return kvs;
}

function KV(key, val) {
    var t = this;
    t.key = ko.observable(key);
    t.val = ko.observable(val);

    t.parent = function() {
        var x = t.key().substring(0,t.key().length - 1);
        var n = x.lastIndexOf("/");
        if (n == -1) {
            return "";
        }
        return t.key().substring(0,n) + "/";
    };

    t.href = ko.computed(function() {
        var x = t.key().replace(location.hash,"#");
        return t.val() === "" ? location.hash + x + "/" : null;
    });

    t.update = function() {
        var key = location.hash + t.key();
        key = key.replace("#","");
        token = Cookies.get("acl_token");
        // make api call here
        $.ajax(url + "/v1/kv" + "?token=" + token, {
            data: ko.toJSON([{ key: key, value: t.val(), token: token }]),
            type: "post", contentType: "application/json",
            success: function() {
                $("#update-message").text("Updated").addClass("update-success").fadeIn().delay(1000).fadeOut();
            },
            error: function() {
                $("#update-message").text("Failed").addClass("update-failure").fadeIn().delay(1000).fadeOut();
            }
        });
    };
}

function KeyValueStoreModel() {
    var t = this;

    t.aclToken = ko.observable(Cookies.get("acl_token") || "");

    // initialize kvs. all else fails if this fails
    t.kvs = ko.observableArray(getAllKVs(t.aclToken()));
    t.refresh = function() {
        t.aclToken(Cookies.get("acl_token"));
        t.kvs(getAllKVs(t.aclToken()));
    };

    t.chosenFolderId = ko.observable("");

    // t.goToFolder = function(folder) {
    //   // location.hash = t.chosenFolderId() + folder.key()
    // };

    t.goBackFromFolder = function() {
        var currentFolder = new KV(t.chosenFolderId(), "");
        location.hash = currentFolder.parent();
        t.chosenFolderId(currentFolder.parent());
    };

    t.folders = ko.computed(function() {
        var seen = new Set();
        return t.kvs().reduce(function (acc, kv) {
            // take keys only from current folder
            if (kv.key().startsWith(t.chosenFolderId())) {
                var k = kv.key().replace(t.chosenFolderId(),"");
                var v = k.indexOf("/") == -1 ? kv.val() : "";

                // display only top folder/key name
                k = k.split("/")[0];

                // take only unique folders
                if (!seen.has(k)) {
                    seen.add(k);
                    acc.push(new KV(k, v));
                }
            }
            return acc;
        }, []);

    });

    t.addKey = ko.observable("");
    t.addVal = ko.observable("");
    t.add = function() {
        var kv = new KV(t.addKey(), t.addVal());
        kv.update();
        t.addKey("");
        t.addVal("");
        t.refresh();
    };

    t.del = function(kv) {
        var key = t.chosenFolderId() + kv.key();
        key = key.replace("#","");
        var recurse = false;
        if (kv.val() === "") {
            key += "/";
            recurse = true;
        }
        // make api call here
        $.ajax(url + "/v1/kv" + "?token=" + t.aclToken(), {
            type: "delete", contentType: "application/json",
            data: ko.toJSON([key]),
            success: function() {
                $("#update-message").text("Deleted").addClass("update-success").fadeIn().delay(1000).fadeOut();
                t.refresh();
            },
            error: function() {
                $("#update-message").text("Failed").addClass("update-failure").fadeIn().delay(1000).fadeOut();
            }
        });
    };

    Sammy(function() {
        this.get('#(.*)', function() {
            t.chosenFolderId(this.params['splat'][0]);
        });

        this.get('/');
    }).run();

}

function AppModel() {
    var t = this;

    t.aclToken = ko.observable(Cookies.get("acl_token") || "");
    t.refresh = function() {
        Cookies.set("acl_token", t.aclToken());
        t.kv.refresh();
    };

    t.kv = new KeyValueStoreModel();
    t.displayKV = ko.observable(false);
    t.displayACL = ko.observable(false);
    t.setKV = function() {
        t.displayKV(true);
        t.displayACL(false);
        location.hash = "#";
        t.kv.chosenFolderId("");
    };
    t.setACL = function() {
        t.displayACL(true);
        t.displayKV(false);
        location.hash = "#";
        t.kv.chosenFolderId("");
    };
}

// Activates knockout.js
ko.applyBindings(new AppModel());
