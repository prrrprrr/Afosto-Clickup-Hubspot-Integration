package clickup

import (
	"encoding/json"
	"log"
	"strings"
	"testing"
)

func TestBlocksToHTML_Image(t *testing.T) {
	raw := `[
        {"text": "oke sick", "attributes": {}},
        {"text": "\n", "attributes": {"block-id": "block-1"}},
        {"type": "image", "text": "image.png",
         "image": {"url": "https://example.com/img.png", "name": "image.png"},
         "attributes": {"width": "300"}}
    ]`
	var blocks []CommentBlock
	if err := json.Unmarshal([]byte(raw), &blocks); err != nil {
		t.Fatal(err)
	}
	got := BlocksToHTML(blocks)
	if !strings.Contains(got, `<img src="https://example.com/img.png"`) {
		t.Errorf("expected img tag, got: %s", got)
	}
	if !strings.Contains(got, `<p>oke sick`) {
		t.Errorf("expected paragraph, got: %s", got)
	}
}

func TestBlocksToHTML_PlainNewline(t *testing.T) {
	raw := `[
        {"text": "line one", "attributes": {}},
        {"text": "\n", "attributes": {"block-id": "block-1"}},
        {"text": "\n", "attributes": {"block-id": "block-2"}},
        {"text": "line two", "attributes": {}}
    ]`
	var blocks []CommentBlock
	if err := json.Unmarshal([]byte(raw), &blocks); err != nil {
		t.Fatal(err)
	}
	got := BlocksToHTML(blocks)
	if strings.Contains(got, "<br>") {
		t.Errorf("should not have <br> for block-id newlines, got: %s", got)
	}
	if !strings.Contains(got, "<p>line one") || !strings.Contains(got, "<p>line two") {
		t.Errorf("expected two paragraphs, got: %s", got)
	}
}
func TestBlocksToHTML_BigTest(t *testing.T) {
	raw := `[
        {
          "text": "Bundles should be first priority: Show bundles on website",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-m4vXsGo31G",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": " - Blotter and pellets are different products",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-vXB4A1AEK7",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": "Surcharge: We need to know what your team needs as an import: What kind of field would need to be send/made for the imported field to be send in the correct place. Surcharge is also a WooCommerce fee: Hook woo commerce cart calculate fees.",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-oHSKZ69jtT",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": "Remote area surcharge and tip charge on the invoice: Ordertip plugin: It is ordered as a fee.",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-HnBDzsLlAh",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": "Sequence of printed labels should be the same as the sequence of the :",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-6zzIQUYvgs",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": " Picklist is now sorted with unique SKU’s, filtered by the position of the location.",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-EuaCmnt8U3",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": "Number of unique orders based on the position in the warehouse.",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-v6Y433YxsI",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": "Paid status is not being set to the grid",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-iAXRUBcMeA",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": "Always release orders from WooCommerce: it is not needed, as in the scan app it also says whether or not the province is missing/needed.",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-UYzFMBhd0l",
            "list": {
              "list": "bullet"
            }
          }
        },
        {
          "text": "If you order 2 bundles, only 1 gets imported.",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-UGtQaoPbij"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-5jqaenyaoO"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-zKVLj5pgQE"
          }
        },
        {
          "text": "Picklist single order lines: it separates orders by",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-z5m2NTeNcB"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-CWAjjw4Dku"
          }
        },
        {
          "text": "Picklist works; but sequence is wrong. Shipping labels are not being sequenced correctly.",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-48NT8V2lsp"
          }
        },
        {
          "type": "link_mention",
          "link_mention": {
            "url": "https://lsd-legal.com/nl/shop/starter-bundle/"
          },
          "text": ""
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-G1C0io26Fs"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-qetP6LIh-a"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-W8p8dkC7GJ"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block--ZmiQJQvhS"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-nXtwVZc1Cz"
          }
        },
        {
          "text": "First all SMP10 orders, then lastly the larger orders with SMP25,",
          "attributes": {}
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-P3yUPYlRd6"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-SPK4I7ucg2"
          }
        },
        {
          "type": "image",
          "text": ")%20Thank%20you%20Tengert.png",
          "image": {
            "id": "895ee705-caa2-41e4-a448-97c9e5332a0e.png",
            "name": ")%20Thank%20you%20Tengert.png",
            "title": ")%20Thank%20you%20Tengert.png",
            "type": "png",
            "extension": "image/png",
            "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png",
            "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png",
            "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png",
            "url": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png",
            "uploaded": true
          },
          "attributes": {
            "width": "300",
            "data-id": "895ee705-caa2-41e4-a448-97c9e5332a0e.png",
            "data-attachment": {
              "id": "895ee705-caa2-41e4-a448-97c9e5332a0e.png",
              "version": "0",
              "date": 1774522346314,
              "name": ") Thank you Tengert.png",
              "title": ") Thank you Tengert.png",
              "extension": "png",
              "source": 1,
              "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e_small.png",
              "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e_medium.png",
              "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png",
              "width": 900,
              "height": 446,
              "url": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png",
              "url_w_query": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png?view=open",
              "url_w_host": "https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png"
            },
            "data-natural-width": "900",
            "data-natural-height": "446"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-OVGIzA60-l"
          }
        },
        {
          "type": "image",
          "text": "Marcin%20Tene.png",
          "image": {
            "id": "b15b33d3-86e1-4641-9b31-5b9f8407d689.png",
            "name": "Marcin%20Tene.png",
            "title": "Marcin%20Tene.png",
            "type": "png",
            "extension": "image/png",
            "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png",
            "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png",
            "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png",
            "url": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png",
            "uploaded": true
          },
          "attributes": {
            "width": "300",
            "data-id": "b15b33d3-86e1-4641-9b31-5b9f8407d689.png",
            "data-attachment": {
              "id": "b15b33d3-86e1-4641-9b31-5b9f8407d689.png",
              "version": "0",
              "date": 1774522351958,
              "name": "Marcin Tene.png",
              "title": "Marcin Tene.png",
              "extension": "png",
              "source": 1,
              "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689_small.png",
              "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689_medium.png",
              "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png",
              "width": 708,
              "height": 583,
              "url": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png",
              "url_w_query": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png?view=open",
              "url_w_host": "https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png"
            },
            "data-natural-width": "708",
            "data-natural-height": "583"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-jIxYA00Mn9"
          }
        },
        {
          "type": "image",
          "text": "Versand.png",
          "image": {
            "id": "53df8e26-e395-4dd1-8c89-e34404b72e8e.png",
            "name": "Versand.png",
            "title": "Versand.png",
            "type": "png",
            "extension": "image/png",
            "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png",
            "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png",
            "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png",
            "url": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png",
            "uploaded": true
          },
          "attributes": {
            "width": "300",
            "data-id": "53df8e26-e395-4dd1-8c89-e34404b72e8e.png",
            "data-attachment": {
              "id": "53df8e26-e395-4dd1-8c89-e34404b72e8e.png",
              "version": "0",
              "date": 1774522364783,
              "name": "Versand.png",
              "title": "Versand.png",
              "extension": "png",
              "source": 1,
              "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e_small.png",
              "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e_medium.png",
              "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png",
              "width": 348,
              "height": 378,
              "url": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png",
              "url_w_query": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png?view=open",
              "url_w_host": "https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png"
            },
            "data-natural-width": "348",
            "data-natural-height": "378"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-jXaY-lsqwT"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-dyPaDhdcGZ"
          }
        },
        {
          "type": "image",
          "text": "Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
          "image": {
            "id": "a3a66b3c-86ce-465c-bc5c-414e5d5537ce.png",
            "name": "Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
            "title": "Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
            "type": "png",
            "extension": "image/png",
            "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
            "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
            "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
            "url": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
            "uploaded": true
          },
          "attributes": {
            "width": "300",
            "data-id": "a3a66b3c-86ce-465c-bc5c-414e5d5537ce.png",
            "data-attachment": {
              "id": "a3a66b3c-86ce-465c-bc5c-414e5d5537ce.png",
              "version": "0",
              "date": 1774522361218,
              "name": "Scherm­afbeelding 2026-03-26 om 11.44.38.png",
              "title": "Scherm­afbeelding 2026-03-26 om 11.44.38.png",
              "extension": "png",
              "source": 1,
              "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce_small.png",
              "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce_medium.png",
              "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce_large.png",
              "width": 1645,
              "height": 736,
              "url": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png",
              "url_w_query": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png?view=open",
              "url_w_host": "https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png"
            },
            "data-natural-width": "1645",
            "data-natural-height": "736"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-spI5GkLLek"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-_YxOVcoVYq"
          }
        },
        {
          "type": "image",
          "text": "%C2%A1iiiiiiiiiil.png",
          "image": {
            "id": "37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0.png",
            "name": "%C2%A1iiiiiiiiiil.png",
            "title": "%C2%A1iiiiiiiiiil.png",
            "type": "png",
            "extension": "image/png",
            "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png",
            "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png",
            "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png",
            "url": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png",
            "uploaded": true
          },
          "attributes": {
            "width": "300",
            "data-id": "37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0.png",
            "data-attachment": {
              "id": "37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0.png",
              "version": "0",
              "date": 1774522356628,
              "name": "¡iiiiiiiiiil.png",
              "title": "¡iiiiiiiiiil.png",
              "extension": "png",
              "source": 1,
              "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0_small.png",
              "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0_medium.png",
              "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png",
              "width": 1546,
              "height": 697,
              "url": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png",
              "url_w_query": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png?view=open",
              "url_w_host": "https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png"
            },
            "data-natural-width": "1546",
            "data-natural-height": "697"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-X2ZfjKWE4r"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-ejLGp2bRVT"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-OwK-O7buHI"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-cVGbBwdz17"
          }
        },
        {
          "type": "image",
          "text": "Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
          "image": {
            "id": "92cf53de-7e8b-43b6-93d7-435ca743ea04.png",
            "name": "Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
            "title": "Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
            "type": "png",
            "extension": "image/png",
            "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
            "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
            "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
            "url": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
            "uploaded": true
          },
          "attributes": {
            "width": "300",
            "data-id": "92cf53de-7e8b-43b6-93d7-435ca743ea04.png",
            "data-attachment": {
              "id": "92cf53de-7e8b-43b6-93d7-435ca743ea04.png",
              "version": "0",
              "date": 1774522342428,
              "name": "Scherm­afbeelding 2026-03-26 om 11.12.48.png",
              "title": "Scherm­afbeelding 2026-03-26 om 11.12.48.png",
              "extension": "png",
              "source": 1,
              "thumbnail_small": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04_small.png",
              "thumbnail_medium": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04_medium.png",
              "thumbnail_large": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
              "width": 1170,
              "height": 539,
              "url": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png",
              "url_w_query": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png?view=open",
              "url_w_host": "https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png"
            },
            "data-natural-width": "1170",
            "data-natural-height": "539"
          }
        },
        {
          "text": "\n",
          "attributes": {
            "block-id": "block-AM8pCJ3D-o"
          }
        }
      ]`
	var blocks []CommentBlock
	if err := json.Unmarshal([]byte(raw), &blocks); err != nil {
		t.Fatal(err)
	}
	test := "<ul><li>Bundles should be first priority: Show bundles on website</li></ul>\n<ul><li> - Blotter and pellets are different products</li></ul>\n<ul><li>Surcharge: We need to know what your team needs as an import: What kind of field would need to be send/made for the imported field to be send in the correct place. Surcharge is also a WooCommerce fee: Hook woo commerce cart calculate fees.</li></ul>\n<ul><li>Remote area surcharge and tip charge on the invoice: Ordertip plugin: It is ordered as a fee.</li></ul>\n<ul><li>Sequence of printed labels should be the same as the sequence of the :</li></ul>\n<ul><li> Picklist is now sorted with unique SKU&#39;s, filtered by the position of the location.</li></ul>\n<ul><li>Number of unique orders based on the position in the warehouse.</li></ul>\n<ul><li>Paid status is not being set to the grid</li></ul>\n<ul><li>Always release orders from WooCommerce: it is not needed, as in the scan app it also says whether or not the province is missing/needed.</li></ul>\n<p>If you order 2 bundles, only 1 gets imported.</p>\n<p>Picklist single order lines: it separates orders by</p>\n<p>Picklist works; but sequence is wrong. Shipping labels are not being sequenced correctly.</p>\n<a href=\"https://lsd-legal.com/nl/shop/starter-bundle/\" target=\"_blank\" rel=\"noopener\">https://lsd-legal.com/nl/shop/starter-bundle/</a>\n<p>First all SMP10 orders, then lastly the larger orders with SMP25,</p>\n<img src=\"https://t1373770.p.clickup-attachments.com/t1373770/895ee705-caa2-41e4-a448-97c9e5332a0e/)%20Thank%20you%20Tengert.png\" alt=\")%20Thank%20you%20Tengert.png\" width=\"300\">\n<img src=\"https://t1373770.p.clickup-attachments.com/t1373770/b15b33d3-86e1-4641-9b31-5b9f8407d689/Marcin%20Tene.png\" alt=\"Marcin%20Tene.png\" width=\"300\">\n<img src=\"https://t1373770.p.clickup-attachments.com/t1373770/53df8e26-e395-4dd1-8c89-e34404b72e8e/Versand.png\" alt=\"Versand.png\" width=\"300\">\n<img src=\"https://t1373770.p.clickup-attachments.com/t1373770/a3a66b3c-86ce-465c-bc5c-414e5d5537ce/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png\" alt=\"Scherm%C2%ADafbeelding%202026-03-26%20om%2011.44.38.png\" width=\"300\">\n<img src=\"https://t1373770.p.clickup-attachments.com/t1373770/37ab0f0b-0ad3-4e4f-8da3-1f144b6990b0/%C2%A1iiiiiiiiiil.png\" alt=\"%C2%A1iiiiiiiiiil.png\" width=\"300\">\n<img src=\"https://t1373770.p.clickup-attachments.com/t1373770/92cf53de-7e8b-43b6-93d7-435ca743ea04/Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png\" alt=\"Scherm%C2%ADafbeelding%202026-03-26%20om%2011.12.48.png\" width=\"300\">"
	got := BlocksToHTML(blocks)
	if strings.Contains(got, test) {
		//t.Errorf("should not have <br> for block-id newlines, got: %s", got)
		t.Errorf("correct")
	}

	log.Print(got)
	log.Print(test)
}
