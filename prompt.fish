# Render background, as this handler has file descriptor
function __render_prompt_bg --on-event fish_prompt
  ./bedotia image
end

# This one renders the text, after the image has been rendered
function fish_prompt 
  set_color -o black
  ./bedotia text
end

