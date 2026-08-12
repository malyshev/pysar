import { CornerTicks } from "@/components/corner-ticks";
import { buttondownSubscribeAction } from "@/lib/site";

/**
 * Boxsi Footer subscribe cell — structure from LandingTemplate Footer.tsx.
 * Native POST to Buttondown embed endpoint (dec-20260808-footer-subscribe-host-form-598699a0).
 * No API key; browser navigates to Buttondown’s confirm / double-opt-in response.
 */
export function FooterSubscribe() {
  return (
    <form
      action={buttondownSubscribeAction}
      method="post"
      className="order-1 flex grow flex-col border-b-2 border-zinc-800 md:order-2 md:border-b-0 md:border-s-2"
    >
      <p className="border-b border-zinc-800 px-4 py-2.5 text-xs text-zinc-400">
        Occasional product notes — no spam.
      </p>
      <div className="flex min-h-0 grow">
        <label className="sr-only" htmlFor="footer-subscribe-email">
          Email for occasional Pysar product notes
        </label>
        <input
          id="footer-subscribe-email"
          name="email"
          type="email"
          autoComplete="email"
          required
          placeholder="Enter your email"
          className="inline-flex h-full w-full items-center justify-center border-0 bg-transparent px-3 py-2.5 text-sm font-medium text-white transition-all duration-500 placeholder:text-zinc-500 focus:outline-none"
        />
        <input type="hidden" name="embed" value="1" />
        <input type="hidden" name="tag" value="getpysar" />
        <div className="relative inline-flex bg-zinc-900 p-4.5">
          <button
            type="submit"
            className="inline-flex cursor-pointer items-center justify-center gap-2 whitespace-nowrap rounded bg-primary px-6 py-2.5 font-medium text-white uppercase transition-all duration-500 hover:bg-primary-hover"
          >
            Subscribe
          </button>
          <CornerTicks className="text-zinc-500" />
        </div>
      </div>
    </form>
  );
}
