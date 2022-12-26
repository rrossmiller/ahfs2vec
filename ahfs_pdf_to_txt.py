from re import T
from PyPDF2 import PdfFileReader
from functools import partial
import os


def extract_information(pdf_path):
    with open(pdf_path, "rb") as f:
        pdf = PdfFileReader(f)
        information = pdf.getDocumentInfo()
        number_of_pages = pdf.getNumPages()

    txt = f"""
    Information about {pdf_path}: 

    Author: {information.author}
    Creator: {information.creator}
    Producer: {information.producer}
    Subject: {information.subject}
    Title: {information.title}
    Number of pages: {number_of_pages}
    """
    # print(txt)

    return information, number_of_pages


def pdf_to_text_all(pdf_path, save_dir, num_pages):
    print()

    f = open(pdf_path, "rb")
    pdf = PdfFileReader(f)
    out = []
    raw = []
    for i in range(num_pages - 1): # skip the last page, which is notes 
        txt = pdf.getPage(i).extract_text()
        raw.append(txt + "\n")
        # for every line
        for l in txt.split("\n"):
            # skip page footers
            if "©" in l:
                l = l.replace(
                    f"© 2019, American Society of Health-System Pharmacists, Inc. Page {i} of 31",
                    "",
                )
            if " ()" in l:
                l = l.replace(" ()", "")

            if l == '':
                continue
            # if ) is in the line and a space does not follow it
            l_split = l.split(")")
            if ")" in l and l_split[1] != "":
                if l_split[1][0] not in [" ", "2"]:
                    # split the line there
                    a = l_split[0] + ")\n"
                    b = l_split[1] + "\n"
                    if a == " ()\n":
                        a = ""
                    if b == " ()\n":
                        b = ""
                    out.extend([a, b])
                else:
                    out.append(l + "\n")

            else:
                if (
                    l
                    == "92:92 Other Miscellaneous Therapeutic AgentsAbobotulinumtoxinA (315012)"
                ):
                    l = "92:92 Other Miscellaneous Therapeutic Agents\nAbobotulinumtoxinA (315012)"
                out.append(l + "\n")
    f.close()

    with open(f"{save_dir}/all_pages.txt", "w") as fout:
        fout.writelines(out)

    with open(f"{save_dir}/raw.txt", "w") as fout:
        fout.writelines(raw)


if __name__ == "__main__":
    save_dir = "AHFSClassificationwithDrugs2019"
    if not os.path.exists(save_dir):
        os.mkdir(save_dir)

    pdf_path = "AHFSClassificationwithDrugs2019.pdf"
    _, num_pages = extract_information(pdf_path)

    # merge pages for test
    pdf_to_text_all(pdf_path, save_dir, num_pages)

    print('there will be some additional manual stuff to do')
    # regex:'\)\d' in the output file and split manually. UNLESS it's F(ab')2
    # regex '© 2019, American Society of Health-System Pharmacists, Inc. Page \d\d of 31\n'  -- remove
    #   and '© 2019, American Society of Health-System Pharmacists, Inc. Page \d of 31\n'    -- remove
    # Aldosterone System \nInhibitors needs to have the new line removed
    # ' \n' ... remove new line
    # other dangling words that are the end of a description (because they have a new line char in the middle)
    # rejoin Norepinephrine-\n
    # split into 3. use pdf for reference: Helidac (bismuth, metronidazole, tetracycline) Pylera (bismuth, metronidazole, tetracycline) 56:32 Prokinetic Agents
    # Aurothioglucose/Gold Sodium Thiomalate split after
    # split House Dust Mites Allergen Extract (Odactra) Short Ragweed Pollen Allergen Extract 80:04 Antitoxins and Immune Globulins
    # Diphtheria and Tetanus Toxoids and Acellular 80:12 Vaccines
