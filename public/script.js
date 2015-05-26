$(document).ready(function() {
  $('a.menu-btn').click(function() {
    $('nav').slideToggle(100);
    return false;
  });

  $(window).resize(function(){
    var w = $(window).width();
    var menu = $('.site-header nav');
    if (w > 680 && menu.is(':hidden')) {
      menu.removeAttr('style');
    }
  });
});
