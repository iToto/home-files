import { Link } from "@remix-run/react";

import { useOptionalUser } from "~/utils";

export default function Index() {
  const user = useOptionalUser();
  return (
    <main className="relative min-h-screen bg-white sm:flex sm:items-center sm:justify-center">
      <div className="relative sm:pb-16 sm:pt-8">
        
        {/* header */}
        <nav className="flex items-center justify-between flex-wrap bg-teal-500 p-6">
          <div className="flex items-center flex-shrink-0 text-white mr-6">
            <img src="https://images.squarespace-cdn.com/content/v1/5b70a3b5365f020c4572a851/1608126715346-PEGQIR50LL1YRAHT5570/Luna+IT-logo+%28small%29.png?format=1500w" alt="Luna IT" className="h-20 w-22 mr-0" />
          </div>
          <div className="block lg:hidden">
            <button className="flex items-center px-3 py-2 border rounded text-teal-200 border-teal-400 hover:text-white hover:border-white">
              <svg className="fill-current h-3 w-3" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><title>Menu</title><path d="M0 3h20v2H0V3zm0 6h20v2H0V9zm0 6h20v2H0v-2z"/></svg>
            </button>
          </div>
          <div className="w-full block flex-grow lg:flex lg:items-center lg:w-auto">
            <div className="text-sm lg:flex-grow">
              <a href="#responsive-header" className="block mt-4 lg:inline-block lg:mt-0 text-teal-200 hover:text-white mr-4">
                Home
              </a>
              <a href="#responsive-header" className="block mt-4 lg:inline-block lg:mt-0 text-teal-200 hover:text-white mr-4">
                Our Team
              </a>
            </div>
            <div>
              <a href="#responsive-header" className="block mt-4 lg:inline-block lg:mt-0 text-teal-200 hover:text-white mr-4">
                Contact Us
              </a>
            </div>
          </div>
        </nav>

        <div className="mx-auto max-w-7xl sm:px-6 lg:px-8">
          <div className="relative shadow-xl sm:overflow-hidden sm:rounded-2xl">
            <div className="absolute inset-0">
              <img
                className="h-full w-full object-cover"
                src="https://images.squarespace-cdn.com/content/v1/5b70a3b5365f020c4572a851/1627781730896-15547C72EUKLRFTIC71U/LunaIT%2B16-9%2BBanner%2BChairs%2BNo%2BLogo.jpg?format=1500w"
                alt="Chairs"
              />
              <div className="absolute inset-0 bg-[color:rgba(254,204,27,0.0)] mix-blend-multiply" />
            </div>

            <div className="relative px-4 pt-16 pb-8 sm:px-6 sm:pt-24 sm:pb-14 lg:px-8 lg:pb-20 lg:pt-32">
              <h1 className="text-center text-6xl font-extrabold tracking-tight sm:text-8xl lg:text-4xl">
                <span className="block text-blue drop-shadow-md">
                  Bespoke Web Development and IT Services
                </span>
              </h1>
              <p className="mx-auto mt-6 max-w-lg text-center text-xl text-green sm:max-w-3xl">
                We are a small team of web developers and IT professionals based in North America. We specialize in bespoke web development and IT services for small businesses and individuals.
              </p>
              {/* <a href="https://remix.run">
                <img
                  src="https://user-images.githubusercontent.com/1500684/158298926-e45dafff-3544-4b69-96d6-d3bcc33fc76a.svg"
                  alt="Remix"
                  className="mx-auto mt-16 w-full max-w-[12rem] md:max-w-[16rem]"
                />
              </a> */}
            </div>
          </div>
        </div>

        
      </div>
    </main>
  );
}
