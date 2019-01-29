'use strict';
/*
window.onload = ev => {
  console.log('yeet');
  var el = document.getElementsByClassName("card");
  console.log(el);
  for(let i = 0; i < el.length; i++) {
    el[i].onclick = e => {
      if (el[i].hasAttribute("data-post")) {
        window.location = "post/" + el[i].getAttribute("data-post");
      }
    };

    var contentDiv = el[i].getElementsByClassName("content")[0];
    console.log(contentDiv);
    if (contentDiv) {
      var content = contentDiv.innerHTML.trim();
      if (content.split(" ").length > 25) {
        content = content.split(" ").slice(0, 25).join(" ") + " ...";
      }
      contentDiv.innerHTML = content;
    }
  }
};*/

$(() => {
  $(".card").each((i, el) => {
    $(el).css("top", "100px").delay(i * 400).animate({
      top: 0
    }, 700, "easeOutQuad");
    $(el).click(e => {
      if (el[i].hasAttribute("data-post")) {
        window.location = "post/" + el[i].getAttribute("data-post");
      }
    });
  })
})
