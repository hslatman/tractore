import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { DateTime } from "luxon";
import React, { FC, useEffect, useState } from "react";
import Client, { monitor, site } from "./client";

const client = new Client(window.location.origin);

function App() {
  return (
    <></>
    // <>
    //   <div className="min-h-full container px-4 mx-auto my-16">
    //     <h2 className="text-2xl font-bold leading-7 text-gray-900 sm:truncate sm:text-3xl sm:tracking-tight">
    //       Uptime Monitoring
    //     </h2>

    //     <main className="pt-8 pb-16">
    //       <SiteList />
    //     </main>
    //   </div>
    // </>
  );
}


export default App;

const validURL = (url: string) => {
  const idx = url.lastIndexOf(".");
  if (idx === -1 || url.substring(idx + 1) === "") {
    return false;
  }

  if (!url.startsWith("http:") && !url.startsWith("https:")) {
    url = "https://" + url;
  }

  try {
    const u = new URL(url);
    return u.protocol === "http:" || u.protocol === "https:";
  } catch (_) {
    return false;
  }
};



const TimeDelta: FC<{ dt: DateTime }> = ({ dt }) => {
  const compute = () => dt.toRelative();
  const [str, setStr] = useState(compute());

  useEffect(() => {
    const handler = () => setStr(compute());
    const timer = setInterval(handler, 1000);
    return () => clearInterval(timer);
  }, [dt]);

  return <>{str}</>;
};