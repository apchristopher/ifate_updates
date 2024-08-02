document.addEventListener('DOMContentLoaded', () => {
    const baseUrl = 'https://www.instituteforapprenticeships.org/apprenticeship-standards/';
    const downloadUrl = `${baseUrl}download`;

    async function fetchData() {
        try {
            // Fetch the main page content
            const response = await fetch(baseUrl);
            const pageText = await response.text();
            const parser = new DOMParser();
            const doc = parser.parseFromString(pageText, 'text/html');

            // Extract data-standard attributes
            const standardIds = Array.from(doc.querySelectorAll('[data-standard]')).map(div => JSON.parse(div.dataset.standard).id);

            // Fetch CSV content
            const csvResponse = await fetch(downloadUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify(standardIds)
            });

            const csvText = await csvResponse.text();
            processCsv(csvText);

        } catch (error) {
            console.error('Error fetching data:', error);
        }
    }

    function processCsv(csvText) {
        // Parse CSV data
        const rows = csvText.split('\n').slice(1); // Skip the header row
        const cutoffDate = new Date();
        cutoffDate.setDate(cutoffDate.getDate() - 14);

        const data = rows.map(row => {
            const cols = row.split(',');
            return {
                reference: cols[0],
                name: cols[1],
                lastUpdated: new Date(cols[2]),
                link: cols[3]
            };
        });

        // Filter and sort the data
        const filteredData = data.filter(item => item.lastUpdated >= cutoffDate)
                                 .sort((a, b) => b.lastUpdated - a.lastUpdated);

        // Format the data for display
        filteredData.forEach(item => {
            item.version = extractVersion(item.link);
            item.lastUpdatedFormatted = item.lastUpdated.toISOString().split('T')[0];
        });

        displayData(filteredData);
    }

    function extractVersion(link) {
        const parts = link.split('-v');
        return parts.length > 1 ? parts[1].replace('-', '.') : '';
    }

    function displayData(data) {
        const container = document.getElementById('data-container');
        container.innerHTML = data.map(item => `
            <div class="item">
                <p><strong>${item.reference} ${item.name} ${item.version}</strong></p>
                <p>Last Updated: ${item.lastUpdatedFormatted}</p>
                <p><a href="${item.link}">Link</a></p>
            </div>
        `).join('');
    }

    fetchData();
});
