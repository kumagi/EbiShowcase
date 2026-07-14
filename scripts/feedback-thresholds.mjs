const count = Number(process.argv[2] || 0);
if (count >= 500) console.log("ARCHIVE: export rows, clear the response sheet, preserve the archive.");
else if (count >= 20) console.log("TRIAGE: review the newest 20 responses and mark each handled/not-planned.");
else console.log("COLLECT: no threshold reached.");
