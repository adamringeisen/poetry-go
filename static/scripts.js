// Display (or hide) into block below poem
function showBar(div_id) {
  const x = document.getElementById(div_id);
  if (x.style.display == "none") {
    x.style.display = "block";
  } else {
    x.style.display = "none";
  }
}
