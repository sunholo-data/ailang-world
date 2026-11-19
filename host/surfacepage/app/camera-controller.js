// The two projections share one document, one graph and a reversible camera.
(() => {
  if (typeof window === 'undefined') return;
  function create(stage, space, hooks) {
    const scene = WorldHeadScene.create(space, hooks);
    let angle = 0, destination = 0, frame = 0, graph, layout, side;
    let selected = '', keyboard = false;
    const paint = () => {
      stage.style.setProperty('--side-visibility', String(Math.max(0, 1 - angle * 2)));
      scene.render(graph, layout, side, angle, selected);
    };
    function update() {
      if (space.hidden) return;
      graph = hooks.model();
      if (!graph) return;
      layout = WorldHeadCamera.layout(graph, {width:stage.clientWidth, height:stage.clientHeight,
        selected, attention:hooks.attention()});
      side = hooks.sidePoints();
      paint();
    }
    function settle() {
      frame = 0;
      angle = destination;
      if (angle === 1) {
        paint();
        space.inert = false;
        stage.dataset.view = 'heads';
        if (keyboard) scene.focus(selected);
      } else {
        space.hidden = true;
        delete stage.dataset.camera;
        stage.dataset.view = graph ? 'timeline' : stage.dataset.view;
        hooks.sideVisible(keyboard);
      }
    }
    function turn(heads, id, instant, focus) {
      if (!hooks.model()) return;
      cancelAnimationFrame(frame);
      destination = heads ? 1 : 0;
      selected = id || selected;
      keyboard = focus;
      if (!heads && space.hidden) return;
      space.hidden = false;
      space.inert = true;
      stage.dataset.camera = 'turn';
      stage.dataset.view = 'turning';
      update();
      if (instant || hooks.reduced() || destination === angle) return settle();
      const start = angle, time = performance.now();
      const span = 1250 * Math.abs(destination - start);
      const step = now => {
        const t = Math.min(1, (now - time) / span);
        if (t === 1 || hooks.reduced()) return settle();
        angle = start + (destination - start) * (t*t*(3 - 2*t));
        paint();
        frame = requestAnimationFrame(step);
      };
      frame = requestAnimationFrame(step);
    }
    function hide() {
      cancelAnimationFrame(frame);
      frame = 0;
      space.hidden = space.inert = true;
      delete stage.dataset.camera;
    }
    return {turn, update, hide,
      select(id) { if (id !== selected && angle === 1 && !space.hidden) { selected = id; paint(); } },
      get active() { return !space.hidden; },
      get heads() { return destination === 1; }
    };
  }
  window.WorldCameraController = {create};
})();
