import { JSX } from 'solid-js';

const TaskDialog = (props: JSX.DialogHtmlAttributes<HTMLDialogElement>) => {
  return (
    <>
      <dialog class="modal" {...props}>
          <div class="modal-box">
            <form method="dialog">
              <input type="text" placeholder="New Task..." class="input input-ghost" />
              <fieldset class="fieldset w-xs bg-base-200 border border-base-300 p-4 rounded-box">
                <legend class="fieldset-legend">Page details</legend>
                
                <label class="fieldset-label">Title</label>
                <input type="text" class="input" placeholder="My awesome page" />
                
                <label class="fieldset-label">Slug</label>
                <input type="text" class="input" placeholder="my-awesome-page" />
                
                <label class="fieldset-label">Author</label>
                <input type="text" class="input" placeholder="Name" />
              </fieldset>
            </form>
            <div class="modal-action">
              <button class="btn">Close</button>
            </div>
          </div>
        </dialog>
    </>
  );
}

export default TaskDialog